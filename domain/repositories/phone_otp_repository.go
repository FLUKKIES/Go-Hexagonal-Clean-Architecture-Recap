package repositories

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type IPhoneOTPRepository interface {
	Create(otp *entities.PhoneOTP) error
	// FindLatestByPhone ใช้เช็ค Rate Limit (60 วินาที) และหา OTP ล่าสุด
	FindLatestByPhone(phoneNumber string) (*entities.PhoneOTP, error)
	MarkAsUsed(id uuid.UUID) error
}
