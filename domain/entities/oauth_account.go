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
	ID         uuid.UUID
	UserID     uuid.UUID
	Provider   OAuthProvider 
	ProviderID string        // User ID จาก Provider
	Email      string        // Email ที่ได้จาก Provider
	CreatedAt  time.Time
}
