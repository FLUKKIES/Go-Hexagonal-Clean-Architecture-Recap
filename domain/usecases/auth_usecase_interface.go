package usecases

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/requests"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/responses"
	"github.com/google/uuid"
)

type IAuthUsecase interface {
	// Email/Password Auth
	Register(req *requests.RegisterRequest) (*responses.LoginResponse, error)
	Login(req *requests.LoginRequest) (*responses.LoginResponse, error)

	// Session Management
	RefreshToken(refreshToken string) (*responses.LoginResponse, error)
	Logout(sessionID uuid.UUID) error
	LogoutAll(userID uuid.UUID) error // Kick ทุก Device

	// OAuth — provider = "google" | "facebook"
	GetOAuthURL(provider, state string) (string, error)
	HandleOAuthCallback(provider, code string) (*responses.LoginResponse, error)

	// Password Reset
	ForgotPassword(email string) error
	ResetPassword(req *requests.ResetPasswordRequest) error

	// Phone OTP Verification
	SendPhoneOTP(req *requests.SendOTPRequest) error
	VerifyPhoneOTP(userID uuid.UUID, req *requests.VerifyOTPRequest) error
}
