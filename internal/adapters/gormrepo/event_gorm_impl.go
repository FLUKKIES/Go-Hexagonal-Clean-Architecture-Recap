package gormrepo

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepositoryImpl(db *gorm.DB) repositories.IEventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) Create(event *entities.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepositoryImpl) FindByID(id uuid.UUID) (*entities.Event, error) {
	var event entities.Event
	if err := r.db.First(&event, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventRepositoryImpl) FindAll() ([]entities.Event, error) {
	var events []entities.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}


