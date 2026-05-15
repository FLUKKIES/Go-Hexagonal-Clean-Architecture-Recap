package model

import (
	"time"

	"github.com/google/uuid"
)

type OAuthAccountGormModel struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	Provider   string    `gorm:"not null"`                                    // "google", "facebook"
	ProviderID string    `gorm:"not null"`                                    // User ID จาก Provider
	Email      string    `gorm:"not null"`
	CreatedAt  time.Time
}

func (OAuthAccountGormModel) TableName() string { return "oauth_accounts" }
