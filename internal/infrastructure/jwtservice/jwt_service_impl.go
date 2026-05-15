package jwtservice

import (
	"errors"
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const accessTokenExpiry = 15 * time.Minute

type jwtServiceImpl struct {
	secretKey []byte
}

func NewJWTService(secretKey string) ports.IJWTService {
	return &jwtServiceImpl{secretKey: []byte(secretKey)}
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *jwtServiceImpl) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := jwtCustomClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtServiceImpl) ValidateAccessToken(tokenStr string) (*ports.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user_id in token")
	}

	return &ports.JWTClaims{
		UserID: userID,
		Role:   claims.Role,
	}, nil
}
