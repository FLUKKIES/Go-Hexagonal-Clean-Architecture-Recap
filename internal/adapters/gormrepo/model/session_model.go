package model

import (
	"time"

	"github.com/google/uuid"
)

type SessionGormModel struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID           uuid.UUID `gorm:"type:uuid;index;not null"`
	RefreshTokenHash string    `gorm:"uniqueIndex;not null"`
	UserAgent        string
	ClientIP         string
	ExpiresAt        time.Time `gorm:"index"`
	CreatedAt        time.Time
}

func (SessionGormModel) TableName() string { return "sessions" }
