package model

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetTokenGormModel struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	TokenHash string     `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"index"`
	UsedAt    *time.Time // nil = ยังไม่ถูกใช้
	CreatedAt time.Time
}

func (PasswordResetTokenGormModel) TableName() string { return "password_reset_tokens" }
