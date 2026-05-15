package rest

import (
	"errors"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/exceptions"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/requests"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/responses"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/usecases"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct {
	usecase    usecases.IAuthUsecase
	jwtService ports.IJWTService
}

func NewAuthController(uc usecases.IAuthUsecase, jwtService ports.IJWTService) *AuthController {
	return &AuthController{usecase: uc, jwtService: jwtService}
}

// POST /api/auth/register
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}

	resp, err := c.usecase.Register(&req)
	if err != nil {
		return c.handleError(ctx, err)
	}
	c.setAuthCookies(ctx, resp)
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// POST /api/auth/login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}
	resp, err := c.usecase.Login(&req)
	if err != nil {
		return c.handleError(ctx, err)
	}
	c.setAuthCookies(ctx, resp)
	return ctx.JSON(resp)
}

// POST /api/auth/refresh
func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	// ดึง Refresh Token จาก Cookie ก่อน (ถ้าไม่มี ค่อยไปหาใน Body เผื่อรองรับแบบเก่า)
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := ctx.BodyParser(&body); err == nil {
			refreshToken = body.RefreshToken
		}
	}

	if refreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token is required"})
	}

	resp, err := c.usecase.RefreshToken(refreshToken)
	if err != nil {
		return c.handleError(ctx, err)
	}
	c.setAuthCookies(ctx, resp)
	return ctx.JSON(resp)
}

// POST /api/auth/logout
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	sessionID, err := c.extractSessionID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := c.usecase.Logout(sessionID); err != nil {
		return c.handleError(ctx, err)
	}
	c.clearAuthCookies(ctx)
	return ctx.JSON(fiber.Map{"message": "logged out successfully"})
}

// POST /api/auth/logout-all
func (c *AuthController) LogoutAll(ctx *fiber.Ctx) error {
	userID, err := c.extractUserID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := c.usecase.LogoutAll(userID); err != nil {
		return c.handleError(ctx, err)
	}
	c.clearAuthCookies(ctx)
	return ctx.JSON(fiber.Map{"message": "logged out from all devices"})
}

// POST /api/auth/forgot-password
func (c *AuthController) ForgotPassword(ctx *fiber.Ctx) error {
	var req requests.ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}
	// ส่ง response เหมือนกันเสมอ — ป้องกัน Email Enumeration Attack
	_ = c.usecase.ForgotPassword(req.Email)
	return ctx.JSON(fiber.Map{"message": "if this email exists, a reset link has been sent"})
}

// POST /api/auth/reset-password
func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	var req requests.ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}
	if err := c.usecase.ResetPassword(&req); err != nil {
		return c.handleError(ctx, err)
	}
	return ctx.JSON(fiber.Map{"message": "password reset successfully"})
}

// POST /api/auth/phone/send-otp
func (c *AuthController) SendPhoneOTP(ctx *fiber.Ctx) error {
	var req requests.SendOTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}
	if err := c.usecase.SendPhoneOTP(&req); err != nil {
		return c.handleError(ctx, err)
	}
	return ctx.JSON(fiber.Map{"message": "OTP sent successfully"})
}

// POST /api/auth/phone/verify-otp
func (c *AuthController) VerifyPhoneOTP(ctx *fiber.Ctx) error {
	userID, err := c.extractUserID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	var req requests.VerifyOTPRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}
	if err := c.usecase.VerifyPhoneOTP(userID, &req); err != nil {
		return c.handleError(ctx, err)
	}
	return ctx.JSON(fiber.Map{"message": "phone number verified successfully"})
}

// ─── Error Mapper ─────────────────────────────────────────────────────────────

func (c *AuthController) handleError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, exceptions.ErrDuplicateEmail):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrInvalidCredentials):
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrWeakPassword):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrPasswordMismatch):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrInvalidToken),
		errors.Is(err, exceptions.ErrTokenExpired),
		errors.Is(err, exceptions.ErrSessionExpired):
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrUserNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrOTPRateLimited):
		return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrInvalidOTP):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrPhoneAlreadyVerified):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

// ─── Token Helpers ────────────────────────────────────────────────────────────

// setAuthCookies ตั้งค่า HTTP-Only Cookies สำหรับ Access และ Refresh Token
func (c *AuthController) setAuthCookies(ctx *fiber.Ctx, resp *responses.LoginResponse) {
	// Access Token Cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		HTTPOnly: true,
		Secure:   true, // ใช้ Secure เสมอ (ใน Production ควรเป็น HTTPS)
		SameSite: "Lax",
		MaxAge:   resp.ExpiresIn, // อายุของ Cookie ตรงกับ Token
	})

	// Refresh Token Cookie (อายุ 30 วัน หรือ 720h)
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   resp.RefreshExpiresIn,
	})
}

// clearAuthCookies ล้าง Cookies ทั้งหมดเมื่อ Logout
func (c *AuthController) clearAuthCookies(ctx *fiber.Ctx) {
	ctx.ClearCookie("access_token")
	ctx.ClearCookie("refresh_token")
}

// extractUserID ดึง UserID จาก Context ที่ AuthMiddleware แนบมาให้
func (c *AuthController) extractUserID(ctx *fiber.Ctx) (uuid.UUID, error) {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, errors.New("unauthorized: user_id not found in context")
	}
	return uuid.Parse(userIDStr)
}

// extractSessionID ดึง SessionID จาก X-Session-ID header
// (Client ต้องส่ง Session ID ที่ได้รับตอน Login มาด้วย)
func (c *AuthController) extractSessionID(ctx *fiber.Ctx) (uuid.UUID, error) {
	raw := ctx.Get("X-Session-ID")
	return uuid.Parse(raw)
}
