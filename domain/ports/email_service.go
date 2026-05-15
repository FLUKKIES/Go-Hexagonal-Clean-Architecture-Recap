package ports

// IEmailService — Port สำหรับส่ง Email
// Domain ไม่รู้ว่าใช้ SMTP, SendGrid, หรือ Resend
type IEmailService interface {
	// ส่ง Reset Password Token Email
	SendPasswordResetEmail(toEmail, firstName, resetLink string) error
	// ส่ง Welcome Email มีไว้สำหรับส่งให้กับผู้ใช้ที่สมัครใหม่
	SendWelcomeEmail(toEmail, firstName string) error
}
