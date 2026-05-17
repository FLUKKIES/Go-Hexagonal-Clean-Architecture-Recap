package repositories

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type IEventParticipantRepository interface {
	Create(participant *entities.EventParticipant) error
	FindByEventAndUser(eventID, userID uuid.UUID) (*entities.EventParticipant, error)
	FindByEventID(eventID uuid.UUID) ([]entities.EventParticipant, error)
	CountByEventID(eventID uuid.UUID) (int64, error)
}
