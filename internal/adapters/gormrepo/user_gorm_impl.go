// ตัวเชื่อมต่อ GORM เข้ากับระบบ — Implement IUserRepository Interface

package gormrepo

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo/model"
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
	m := toUserGormModel(user)
	return r.db.Create(m).Error
}

func (r *userRepositoryImpl) FindByID(id uuid.UUID) (*entities.User, error) {
	m := new(model.UserGormModel)
	if err := r.db.First(m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toUserEntity(m), nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entities.User, error) {
	m := new(model.UserGormModel)
	if err := r.db.Where("email = ?", email).First(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toUserEntity(m), nil
}

func (r *userRepositoryImpl) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	return r.db.Model(&model.UserGormModel{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		}).Error
}

func (r *userRepositoryImpl) UpdateVerifiedAt(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.UserGormModel{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"verified_at": now,
			"updated_at":  now,
		}).Error
}

// ─── Mapper Helpers ───────────────────────────────────────────────────────────

func toUserGormModel(u *entities.User) *model.UserGormModel {
	return &model.UserGormModel{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Password:    u.Password,
		PhoneNumber: u.PhoneNumber,
		ProfileUrl:  u.ProfileUrl,
		Role:        string(u.Role),
		VerifiedAt:  u.VerifiedAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func toUserEntity(m *model.UserGormModel) *entities.User {
	return &entities.User{
		ID:          m.ID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		Email:       m.Email,
		Password:    m.Password,
		PhoneNumber: m.PhoneNumber,
		ProfileUrl:  m.ProfileUrl,
		Role:        entities.UserRole(m.Role),
		VerifiedAt:  m.VerifiedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
