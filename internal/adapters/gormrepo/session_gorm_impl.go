package gormrepo

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepositoryImpl(db *gorm.DB) repositories.ISessionRepository {
	return &sessionRepositoryImpl{db: db}
}

func (r *sessionRepositoryImpl) Create(session *entities.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepositoryImpl) FindByRefreshTokenHash(hash string) (*entities.Session, error) {
	m := new(entities.Session)
	if err := r.db.Where("refresh_token_hash = ?", hash).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *sessionRepositoryImpl) DeleteByID(id uuid.UUID) error {
	return r.db.Delete(&entities.Session{}, "id = ?", id).Error
}

func (r *sessionRepositoryImpl) DeleteAllByUserID(userID uuid.UUID) error {
	return r.db.Delete(&entities.Session{}, "user_id = ?", userID).Error
}
