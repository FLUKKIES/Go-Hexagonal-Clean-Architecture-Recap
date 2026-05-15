package gormrepo

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type phoneOTPRepositoryImpl struct {
	db *gorm.DB
}

func NewPhoneOTPRepositoryImpl(db *gorm.DB) repositories.IPhoneOTPRepository {
	return &phoneOTPRepositoryImpl{db: db}
}

func (r *phoneOTPRepositoryImpl) Create(otp *entities.PhoneOTP) error {
	m := &model.PhoneOTPGormModel{
		ID:          otp.ID,
		PhoneNumber: otp.PhoneNumber,
		OTPHash:     otp.OTPHash,
		ExpiresAt:   otp.ExpiresAt,
		CreatedAt:   otp.CreatedAt,
	}
	return r.db.Create(m).Error
}

// FindLatestByPhone หา OTP ล่าสุดของเบอร์นี้ (ใช้เช็ค Rate Limit และ Verify)
func (r *phoneOTPRepositoryImpl) FindLatestByPhone(phoneNumber string) (*entities.PhoneOTP, error) {
	m := new(model.PhoneOTPGormModel)
	if err := r.db.Where("phone_number = ?", phoneNumber).
		Order("created_at DESC").
		First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entities.PhoneOTP{
		ID:          m.ID,
		PhoneNumber: m.PhoneNumber,
		OTPHash:     m.OTPHash,
		ExpiresAt:   m.ExpiresAt,
		UsedAt:      m.UsedAt,
		CreatedAt:   m.CreatedAt,
	}, nil
}

func (r *phoneOTPRepositoryImpl) MarkAsUsed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.PhoneOTPGormModel{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}
