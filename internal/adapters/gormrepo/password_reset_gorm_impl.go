package gormrepo

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
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
	return r.db.Create(token).Error
}

func (r *passwordResetRepositoryImpl) FindValidByTokenHash(hash string) (*entities.PasswordResetToken, error) {
	m := new(entities.PasswordResetToken)
	// ดึงเฉพาะ Token ที่ยังไม่หมดอายุและยังไม่ถูกใช้
	if err := r.db.Where("token_hash = ? AND expires_at > ? AND used_at IS NULL", hash, time.Now()).
		First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *passwordResetRepositoryImpl) MarkAsUsed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entities.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *passwordResetRepositoryImpl) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&entities.PasswordResetToken{}, "user_id = ?", userID).Error
}
