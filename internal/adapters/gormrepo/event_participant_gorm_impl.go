package gormrepo

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventParticipantRepositoryImpl struct {
	db *gorm.DB
}

func NewEventParticipantRepositoryImpl(db *gorm.DB) repositories.IEventParticipantRepository {
	return &eventParticipantRepositoryImpl{db: db}
}

func (r *eventParticipantRepositoryImpl) Create(participant *entities.EventParticipant) error {
	return r.db.Create(participant).Error
}

func (r *eventParticipantRepositoryImpl) FindByEventAndUser(eventID, userID uuid.UUID) (*entities.EventParticipant, error) {
	var participant entities.EventParticipant
	if err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &participant, nil
}

func (r *eventParticipantRepositoryImpl) FindByEventID(eventID uuid.UUID) ([]entities.EventParticipant, error) {
	var participants []entities.EventParticipant
	err := r.db.Preload("User").Where("event_id = ?", eventID).Find(&participants).Error
	return participants, err
}

func (r *eventParticipantRepositoryImpl) CountByEventID(eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entities.EventParticipant{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}
