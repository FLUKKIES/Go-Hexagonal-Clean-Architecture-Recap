package gormrepo

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type passwordResetRepositoryImpl struct {
	db *gorm.DB
}

func NewPasswordResetRepositoryImpl(db *gorm.DB) repositories.IPasswordResetRepository {
	return &passwordResetRepositoryImpl{db: db}
}

func (r *passwordResetRepositoryImpl) Create(token *entities.PasswordResetToken) error {
	m := &model.PasswordResetTokenGormModel{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	}
	return r.db.Create(m).Error
}

func (r *passwordResetRepositoryImpl) FindValidByTokenHash(hash string) (*entities.PasswordResetToken, error) {
	m := new(model.PasswordResetTokenGormModel)
	// ดึงเฉพาะ Token ที่ยังไม่หมดอายุและยังไม่ถูกใช้
	if err := r.db.Where("token_hash = ? AND expires_at > ? AND used_at IS NULL", hash, time.Now()).
		First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.PasswordResetToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (r *passwordResetRepositoryImpl) MarkAsUsed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.PasswordResetTokenGormModel{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *passwordResetRepositoryImpl) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&model.PasswordResetTokenGormModel{}, "user_id = ?", userID).Error
}
