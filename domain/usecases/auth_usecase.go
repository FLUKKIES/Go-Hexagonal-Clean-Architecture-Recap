package usecases

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/exceptions"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/requests"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/responses"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthConfig struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ResetTokenExpiry   time.Duration
	OTPExpiry          time.Duration
	OTPRateLimit       time.Duration
	AppBaseURL         string
}

type authUsecaseImpl struct {
	userRepo         repositories.IUserRepository
	sessionRepo      repositories.ISessionRepository
	oauthRepo        repositories.IOAuthRepository
	resetTokenRepo   repositories.IPasswordResetRepository
	phoneOTPRepo     repositories.IPhoneOTPRepository
	jwtService       ports.IJWTService
	emailService     ports.IEmailService
	smsService       ports.ISMSService
	oauthProviders   map[string]ports.IOAuthProvider // "google" -> GoogleProvider, "facebook" -> FacebookProvider
	config           AuthConfig
}

func NewAuthUsecase(
	userRepo repositories.IUserRepository,
	sessionRepo repositories.ISessionRepository,
	oauthRepo repositories.IOAuthRepository,
	resetTokenRepo repositories.IPasswordResetRepository,
	phoneOTPRepo repositories.IPhoneOTPRepository,
	jwtService ports.IJWTService,
	emailService ports.IEmailService,
	smsService ports.ISMSService,
	oauthProviders []ports.IOAuthProvider,
	config AuthConfig,
) IAuthUsecase {
	providers := make(map[string]ports.IOAuthProvider)
	for _, p := range oauthProviders {
		providers[p.ProviderName()] = p
	}
	return &authUsecaseImpl{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		oauthRepo:      oauthRepo,
		resetTokenRepo: resetTokenRepo,
		phoneOTPRepo:   phoneOTPRepo,
		jwtService:     jwtService,
		emailService:   emailService,
		smsService:     smsService,
		oauthProviders: providers,
		config:         config,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Register
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) Register(req *requests.RegisterRequest) (*responses.LoginResponse, error) {
	// 1. ตรวจสอบ Password ขั้นต่ำ
	if len(req.Password) < 8 {
		return nil, exceptions.ErrWeakPassword
	}

	// 2. ตรวจสอบ Email ซ้ำ
	existing, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, exceptions.ErrDuplicateEmail
	}

	// 3. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedStr := string(hashed)

	// 4. สร้าง User
	user := &entities.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  &hashedStr,
		Role:      entities.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 5. ส่ง Welcome Email (non-blocking — ไม่ล้มเหลวถ้า email ส่งไม่ได้)
	_ = u.emailService.SendWelcomeEmail(user.Email, user.FirstName)

	// 6. สร้าง Session
	return u.createSession(user, "", "")
}

// ─────────────────────────────────────────────────────────────────────────────
// Login
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) Login(req *requests.LoginRequest) (*responses.LoginResponse, error) {
	// 1. หา User ด้วย Email
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	// ใช้ ErrInvalidCredentials แทน ErrUserNotFound เพื่อป้องกัน User Enumeration Attack
	if user == nil || user.Password == nil {
		return nil, exceptions.ErrInvalidCredentials
	}

	// 2. ตรวจสอบ Password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		return nil, exceptions.ErrInvalidCredentials
	}

	// 3. สร้าง Session
	return u.createSession(user, "", "")
}

// ─────────────────────────────────────────────────────────────────────────────
// RefreshToken — Refresh Token Rotation
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) RefreshToken(refreshToken string) (*responses.LoginResponse, error) {
	// 1. Hash Refresh Token แล้วหา Session
	tokenHash := hashSHA256(refreshToken)
	session, err := u.sessionRepo.FindByRefreshTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, exceptions.ErrInvalidToken
	}

	// 2. ตรวจสอบ Session ไม่หมดอายุ
	if time.Now().After(session.ExpiresAt) {
		_ = u.sessionRepo.DeleteByID(session.ID)
		return nil, exceptions.ErrSessionExpired
	}

	// 3. ดึงข้อมูล User
	user, err := u.userRepo.FindByID(session.UserID)
	if err != nil || user == nil {
		return nil, exceptions.ErrUserNotFound
	}

	// 4. Revoke Session เก่า (Token Rotation — ป้องกัน Replay Attack)
	if err := u.sessionRepo.DeleteByID(session.ID); err != nil {
		return nil, err
	}

	// 5. ออก Session ใหม่
	return u.createSession(user, session.UserAgent, session.ClientIP)
}

// ─────────────────────────────────────────────────────────────────────────────
// Logout
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) Logout(sessionID uuid.UUID) error {
	return u.sessionRepo.DeleteByID(sessionID)
}

func (u *authUsecaseImpl) LogoutAll(userID uuid.UUID) error {
	return u.sessionRepo.DeleteAllByUserID(userID)
}

// ─────────────────────────────────────────────────────────────────────────────
// OAuth
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) GetOAuthURL(provider, state string) (string, error) {
	p, ok := u.oauthProviders[provider]
	if !ok {
		return "", fmt.Errorf("unsupported oauth provider: %s", provider)
	}
	return p.GetAuthURL(state), nil
}

func (u *authUsecaseImpl) HandleOAuthCallback(provider, code string) (*responses.LoginResponse, error) {
	p, ok := u.oauthProviders[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported oauth provider: %s", provider)
	}

	// 1. แลก Code เป็น Profile
	profile, err := p.GetUserProfile(code)
	if err != nil {
		return nil, err
	}

	// 2. หา OAuth Account ที่มีอยู่แล้ว
	oauthAcc, err := u.oauthRepo.FindByProviderAndID(provider, profile.ProviderID)
	if err != nil {
		return nil, err
	}

	var user *entities.User

	if oauthAcc != nil {
		// 3a. มี OAuth Account แล้ว — ดึง User
		user, err = u.userRepo.FindByID(oauthAcc.UserID)
		if err != nil || user == nil {
			return nil, exceptions.ErrUserNotFound
		}
	} else {
		// 3b. ยังไม่มี — ตรวจสอบว่า Email มีในระบบแล้วหรือยัง
		user, err = u.userRepo.FindByEmail(profile.Email)
		if err != nil {
			return nil, err
		}

		if user == nil {
			// 4a. สร้าง User ใหม่ (ไม่มี Password)
			profileUrl := profile.ProfileUrl
			user = &entities.User{
				ID:         uuid.New(),
				FirstName:  profile.FirstName,
				LastName:   profile.LastName,
				Email:      profile.Email,
				Password:   nil, // OAuth user ไม่มี Password
				ProfileUrl: &profileUrl,
				Role:       entities.UserRoleUser,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := u.userRepo.Create(user); err != nil {
				return nil, err
			}
		}
		// 4b. Link OAuth Account เข้ากับ User (ไม่ว่าจะสร้างใหม่หรือมีอยู่แล้ว)
		newOAuthAcc := &entities.OAuthAccount{
			ID:         uuid.New(),
			UserID:     user.ID,
			Provider:   entities.OAuthProvider(provider),
			ProviderID: profile.ProviderID,
			Email:      profile.Email,
			CreatedAt:  time.Now(),
		}
		if err := u.oauthRepo.Create(newOAuthAcc); err != nil {
			return nil, err
		}
	}

	return u.createSession(user, "", "")
}

// ─────────────────────────────────────────────────────────────────────────────
// Password Reset
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) ForgotPassword(email string) error {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	// ถ้าไม่เจอ Email ก็ไม่ Error — ป้องกัน User Enumeration Attack
	if user == nil {
		return nil
	}

	// 1. ยกเลิก Token เก่าทั้งหมด
	_ = u.resetTokenRepo.DeleteByUserID(user.ID)

	// 2. สร้าง Cryptographically Secure Token (32 bytes)
	rawToken, err := generateSecureToken(32)
	if err != nil {
		return err
	}

	// 3. Hash Token ก่อนเก็บลง DB
	tokenHash := hashSHA256(rawToken)

	// 3. กำหนด ExpiresAt ตาม config
	resetToken := &entities.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(u.config.ResetTokenExpiry),
		CreatedAt: time.Now(),
	}
	if err := u.resetTokenRepo.Create(resetToken); err != nil {
		return err
	}

	// 4. ส่ง Email พร้อม Raw Token (ไม่ใช่ Hash)
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", u.config.AppBaseURL, rawToken)
	return u.emailService.SendPasswordResetEmail(user.Email, user.FirstName, resetLink)
}

func (u *authUsecaseImpl) ResetPassword(req *requests.ResetPasswordRequest) error {
	// 1. ตรวจสอบ Password
	if len(req.NewPassword) < 8 {
		return exceptions.ErrWeakPassword
	}

	// 2. Hash Token แล้วหาใน DB
	tokenHash := hashSHA256(req.Token)
	resetToken, err := u.resetTokenRepo.FindValidByTokenHash(tokenHash)
	if err != nil {
		return err
	}
	if resetToken == nil {
		return exceptions.ErrInvalidToken
	}

	// 3. ตรวจสอบ Token ไม่หมดอายุและยังไม่ถูกใช้
	if time.Now().After(resetToken.ExpiresAt) || resetToken.UsedAt != nil {
		return exceptions.ErrTokenExpired
	}

	// 4. Hash Password ใหม่
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 5. อัปเดต Password + Mark Token ว่าใช้แล้ว (atomic-ish)
	if err := u.userRepo.UpdatePassword(resetToken.UserID, string(hashed)); err != nil {
		return err
	}
	if err := u.resetTokenRepo.MarkAsUsed(resetToken.ID); err != nil {
		return err
	}

	// 6. Invalidate Session ทั้งหมด (บังคับ Login ใหม่ทุก Device)
	_ = u.sessionRepo.DeleteAllByUserID(resetToken.UserID)

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Phone OTP
// ─────────────────────────────────────────────────────────────────────────────

func (u *authUsecaseImpl) SendPhoneOTP(req *requests.SendOTPRequest) error {
	// 2. ตรวจสอบ Rate Limit ตาม config
	latestOTP, _ := u.phoneOTPRepo.FindLatestByPhone(req.PhoneNumber)
	if latestOTP != nil && time.Since(latestOTP.CreatedAt) < u.config.OTPRateLimit {
		return exceptions.ErrOTPRateLimited
	}

	// 2. สร้าง OTP 6 หลัก แบบ Cryptographically Secure
	otp, err := generateNumericOTP(6)
	if err != nil {
		return err
	}

	// 3. Hash OTP ก่อนเก็บลง DB
	otpHash := hashSHA256(otp)

	// 4. บันทึก OTP Hash ลง DB (อายุตาม config)
	newOTP := &entities.PhoneOTP{
		ID:          uuid.New(),
		PhoneNumber: req.PhoneNumber,
		OTPHash:     otpHash,
		ExpiresAt:   time.Now().Add(u.config.OTPExpiry),
		CreatedAt:   time.Now(),
	}
	if err := u.phoneOTPRepo.Create(newOTP); err != nil {
		return err
	}

	// 4. ส่ง SMS พร้อม OTP จริง
	return u.smsService.SendOTP(req.PhoneNumber, otp)
}

func (u *authUsecaseImpl) VerifyPhoneOTP(userID uuid.UUID, req *requests.VerifyOTPRequest) error {
	// 1. หา User และตรวจสอบว่า Verify แล้วหรือยัง
	user, err := u.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return exceptions.ErrUserNotFound
	}
	if user.VerifiedAt != nil {
		return exceptions.ErrPhoneAlreadyVerified
	}

	// 2. Hash OTP ที่ส่งมาแล้วเปรียบเทียบกับใน DB
	otpHash := hashSHA256(req.OTP)
	latest, err := u.phoneOTPRepo.FindLatestByPhone(req.PhoneNumber)
	if err != nil {
		return err
	}
	if latest == nil || latest.OTPHash != otpHash {
		return exceptions.ErrInvalidOTP
	}

	// 3. ตรวจสอบ OTP ไม่หมดอายุและยังไม่ถูกใช้
	if time.Now().After(latest.ExpiresAt) || latest.UsedAt != nil {
		return exceptions.ErrInvalidOTP
	}

	// 4. Mark OTP ว่าใช้แล้ว + Update VerifiedAt
	if err := u.phoneOTPRepo.MarkAsUsed(latest.ID); err != nil {
		return err
	}
	return u.userRepo.UpdateVerifiedAt(userID)
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper Functions
// ─────────────────────────────────────────────────────────────────────────────

// createSession สร้าง JWT Access Token + Refresh Token + Session ใน DB
func (u *authUsecaseImpl) createSession(user *entities.User, userAgent, clientIP string) (*responses.LoginResponse, error) {
	// สร้าง Access Token (JWT)
	accessToken, err := u.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// สร้าง Refresh Token (UUID random)
	rawRefreshToken := uuid.New().String() + uuid.New().String() // 72 chars
	refreshTokenHash := hashSHA256(rawRefreshToken)

	session := &entities.Session{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		UserAgent:        userAgent,
		ClientIP:         clientIP,
		ExpiresAt:        time.Now().Add(u.config.RefreshTokenExpiry),
		CreatedAt:        time.Now(),
	}
	if err := u.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return &responses.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     rawRefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int(u.config.AccessTokenExpiry.Seconds()),
		RefreshExpiresIn: int(u.config.RefreshTokenExpiry.Seconds()),
		User:             toUserResponse(user),
	}, nil
}

// hashSHA256 แปลง string เป็น SHA-256 hex string
func hashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// generateSecureToken สร้าง random hex token ขนาด n bytes
func generateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generateNumericOTP สร้าง OTP ตัวเลข n หลัก แบบ Cryptographically Secure
func generateNumericOTP(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", digits, n), nil
}

// toUserResponse แปลง User Entity เป็น UserResponse DTO
func toUserResponse(user *entities.User) responses.UserResponse {
	return responses.UserResponse{
		ID:          user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ProfileUrl:  user.ProfileUrl,
		Role:        string(user.Role),
		IsVerified:  user.VerifiedAt != nil,
	}
}
