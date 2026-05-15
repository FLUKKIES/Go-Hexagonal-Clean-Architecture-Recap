// ตัวเชื่อมต่อ GORM เข้ากับระบบ — Implement IUserRepository Interface

package gormrepo

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) repositories.IUserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) FindByID(id uuid.UUID) (*entities.User, error) {
	m := new(entities.User)
	if err := r.db.First(m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entities.User, error) {
	m := new(entities.User)
	if err := r.db.Where("email = ?", email).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

func (r *userRepositoryImpl) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	return r.db.Model(&entities.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		}).Error
}

func (r *userRepositoryImpl) UpdateVerifiedAt(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entities.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"verified_at": now,
			"updated_at":  now,
		}).Error
}

