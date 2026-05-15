package entities

import (
	"time"

	"github.com/google/uuid"
)

// PhoneOTP ใช้สำหรับยืนยันเบอร์มือถือผ่าน SMS
// เก็บ OTP เป็น Hash เพื่อความปลอดภัย — OTP จริงส่งทาง SMS เท่านั้น
type PhoneOTP struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid"`
	PhoneNumber string     `gorm:"index;not null"`
	RefCode     string     // สำหรับ Ref Code Verification
	OTPHash     string     `gorm:"not null"` // เก็บ SHA-256 hash ของ OTP 6 หลัก
	ExpiresAt   time.Time  `gorm:"index"` // หมดอายุใน 5 นาที
	UsedAt      *time.Time // nil = ยังไม่ถูกใช้ (ใช้ได้ครั้งเดียว)
	CreatedAt   time.Time  `gorm:"index"` // ใช้เช็ค Rate Limit (60 วินาที)
}

func (PhoneOTP) TableName() string { return "phone_otps" }
