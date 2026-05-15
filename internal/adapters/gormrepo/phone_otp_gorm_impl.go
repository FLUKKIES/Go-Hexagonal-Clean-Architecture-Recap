package gormrepo

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
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
	return r.db.Create(otp).Error
}

// FindLatestByPhone หา OTP ล่าสุดของเบอร์นี้ (ใช้เช็ค Rate Limit และ Verify)
func (r *phoneOTPRepositoryImpl) FindLatestByPhone(phoneNumber string) (*entities.PhoneOTP, error) {
	m := new(entities.PhoneOTP)
	if err := r.db.Where("phone_number = ?", phoneNumber).
		Order("created_at DESC").
		First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *phoneOTPRepositoryImpl) MarkAsUsed(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entities.PhoneOTP{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}
