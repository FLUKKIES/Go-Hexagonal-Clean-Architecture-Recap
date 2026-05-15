package requests

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3"`
	LastName  string `json:"last_name" validate:"required,min=3"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type SendOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"` // Format: +66812345678
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	OTP         string `json:"otp" validate:"required,len=6,numeric"` // 6 หลัก
}
