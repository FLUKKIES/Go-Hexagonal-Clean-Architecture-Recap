package model

import (
	"time"

	"github.com/google/uuid"
)

type PhoneOTPGormModel struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid"`
	PhoneNumber string     `gorm:"index;not null"`
	OTPHash     string     `gorm:"not null"`
	ExpiresAt   time.Time  `gorm:"index"`
	UsedAt      *time.Time // nil = ยังไม่ถูกใช้
	CreatedAt   time.Time  `gorm:"index"` // ใช้เช็ค Rate Limit
}

func (PhoneOTPGormModel) TableName() string { return "phone_otps" }
