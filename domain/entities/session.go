package entities

import (
	"time"

	"github.com/google/uuid"
)

// Session เก็บ Refresh Token ของแต่ละอุปกรณ์ (1 User มีได้หลาย Session)
type Session struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID           uuid.UUID `gorm:"type:uuid;index;not null"`
	RefreshTokenHash string    `gorm:"uniqueIndex;not null"` // เก็บ SHA-256 hash ของ Refresh Token จริง
	UserAgent        string    // เช่น "Chrome on Windows"
	ClientIP         string    // เก็บ IP ของ Client
	ExpiresAt        time.Time `gorm:"index"`
	CreatedAt        time.Time
}

func (Session) TableName() string { return "sessions" }
