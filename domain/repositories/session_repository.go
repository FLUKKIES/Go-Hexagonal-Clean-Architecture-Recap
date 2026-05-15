package repositories

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type ISessionRepository interface {
	Create(session *entities.Session) error
	FindByRefreshTokenHash(hash string) (*entities.Session, error)
	DeleteByID(id uuid.UUID) error
	DeleteAllByUserID(userID uuid.UUID) error // Logout จากทุก Device
}
