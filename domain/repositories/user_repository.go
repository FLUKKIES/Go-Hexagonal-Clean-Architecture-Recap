package repositories

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(user *entities.User) error
	FindByID(id uuid.UUID) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	UpdatePassword(userID uuid.UUID, hashedPassword string) error
	UpdateVerifiedAt(userID uuid.UUID) error
}
