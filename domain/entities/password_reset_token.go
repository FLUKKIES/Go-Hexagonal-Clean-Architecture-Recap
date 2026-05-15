package entities

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken ใช้สำหรับ Reset Password ผ่าน Email Link
// เก็บเป็น Hash เพื่อความปลอดภัย — Token จริงส่งทาง Email เท่านั้น
type PasswordResetToken struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	TokenHash string     `gorm:"uniqueIndex;not null"` // เก็บ SHA-256 hash ของ Token จริง
	ExpiresAt time.Time  `gorm:"index"` // หมดอายุใน 15 นาที
	UsedAt    *time.Time // nil = ยังไม่ถูกใช้ (ใช้ได้ครั้งเดียว)
	CreatedAt time.Time
}

func (PasswordResetToken) TableName() string { return "password_reset_tokens" }
