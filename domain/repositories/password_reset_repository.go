package repositories

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type IPasswordResetRepository interface {
	Create(token *entities.PasswordResetToken) error
	FindValidByTokenHash(hash string) (*entities.PasswordResetToken, error)
	MarkAsUsed(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error // ยกเลิก Token เก่าทั้งหมดเมื่อขอใหม่
}
