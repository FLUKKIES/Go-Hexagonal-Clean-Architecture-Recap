package gormrepo

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo/model"
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
	m := &model.SessionGormModel{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		UserAgent:        session.UserAgent,
		ClientIP:         session.ClientIP,
		ExpiresAt:        session.ExpiresAt,
		CreatedAt:        session.CreatedAt,
	}
	return r.db.Create(m).Error
}

func (r *sessionRepositoryImpl) FindByRefreshTokenHash(hash string) (*entities.Session, error) {
	m := new(model.SessionGormModel)
	if err := r.db.Where("refresh_token_hash = ?", hash).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.Session{
		ID:               m.ID,
		UserID:           m.UserID,
		RefreshTokenHash: m.RefreshTokenHash,
		UserAgent:        m.UserAgent,
		ClientIP:         m.ClientIP,
		ExpiresAt:        m.ExpiresAt,
		CreatedAt:        m.CreatedAt,
	}, nil
}

func (r *sessionRepositoryImpl) DeleteByID(id uuid.UUID) error {
	return r.db.Delete(&model.SessionGormModel{}, "id = ?", id).Error
}

func (r *sessionRepositoryImpl) DeleteAllByUserID(userID uuid.UUID) error {
	return r.db.Delete(&model.SessionGormModel{}, "user_id = ?", userID).Error
}
