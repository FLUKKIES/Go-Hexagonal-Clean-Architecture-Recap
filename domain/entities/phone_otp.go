package entities

import (
	"time"

	"github.com/google/uuid"
)

// PhoneOTP ใช้สำหรับยืนยันเบอร์มือถือผ่าน SMS
// เก็บ OTP เป็น Hash เพื่อความปลอดภัย — OTP จริงส่งทาง SMS เท่านั้น
type PhoneOTP struct {
	ID          uuid.UUID
	PhoneNumber string
	RefCode     string     // สำหรับ Ref Code Verification
	OTPHash     string     // เก็บ SHA-256 hash ของ OTP 6 หลัก
	ExpiresAt   time.Time  // หมดอายุใน 5 นาที
	UsedAt      *time.Time // nil = ยังไม่ถูกใช้ (ใช้ได้ครั้งเดียว)
	CreatedAt   time.Time  // ใช้เช็ค Rate Limit (60 วินาที)
}
