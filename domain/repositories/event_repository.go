package repositories

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type IEventRepository interface {
	Create(event *entities.Event) error
	FindByID(id uuid.UUID) (*entities.Event, error)
	FindAll() ([]entities.Event, error)
}


