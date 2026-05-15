package entities

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken ใช้สำหรับ Reset Password ผ่าน Email Link
// เก็บเป็น Hash เพื่อความปลอดภัย — Token จริงส่งทาง Email เท่านั้น
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string     // เก็บ SHA-256 hash ของ Token จริง
	ExpiresAt time.Time  // หมดอายุใน 15 นาที
	UsedAt    *time.Time // nil = ยังไม่ถูกใช้ (ใช้ได้ครั้งเดียว)
	CreatedAt time.Time
}
