package ports

// ISMSService — Port สำหรับส่ง SMS OTP
// Domain ไม่รู้ว่าใช้ Twilio, AWS SNS, หรือ provider อื่น
type ISMSService interface {
	SendOTP(phoneNumber, otp string) error
}
