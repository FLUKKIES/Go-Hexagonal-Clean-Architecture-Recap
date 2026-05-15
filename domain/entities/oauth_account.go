package entities

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider string

const (
	OAuthProviderGoogle   OAuthProvider = "google"
	OAuthProviderFacebook OAuthProvider = "facebook"
)

// OAuthAccount เชื่อม User กับ OAuth Provider (1 User มีได้หลาย Provider)
type OAuthAccount struct {
	ID         uuid.UUID     `gorm:"primaryKey;type:uuid"`
	UserID     uuid.UUID     `gorm:"type:uuid;index;not null"`
	Provider   OAuthProvider `gorm:"not null"` // "google", "facebook"
	ProviderID string        `gorm:"not null"` // User ID จาก Provider
	Email      string        `gorm:"not null"` // Email ที่ได้จาก Provider
	CreatedAt  time.Time
}

func (OAuthAccount) TableName() string { return "oauth_accounts" }
