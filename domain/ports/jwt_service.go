package ports

import "github.com/google/uuid"

type JWTClaims struct {
	UserID uuid.UUID
	Role   string
}

type IJWTService interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	ValidateAccessToken(token string) (*JWTClaims, error)
}
