package exceptions

import "errors"

var (
	ErrOTPRateLimited       = errors.New("please wait 60 seconds before requesting a new OTP")
	ErrInvalidOTP           = errors.New("invalid or expired OTP")
	ErrPhoneAlreadyVerified = errors.New("phone number is already verified")
)
