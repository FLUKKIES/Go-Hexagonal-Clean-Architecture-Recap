package model

import (
	"time"

	"github.com/google/uuid"
)

// ตารางใน Database จริงๆ — มี GORM tags อยู่ที่นี่เท่านั้น
type UserGormModel struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid"`
	FirstName   string     `gorm:"not null"`
	LastName    string     `gorm:"not null"`
	Email       string     `gorm:"uniqueIndex;not null"`
	Password    *string    // nil สำหรับ OAuth users
	PhoneNumber *string    `gorm:"uniqueIndex"` // unique + nullable
	ProfileUrl  *string
	Role        string     `gorm:"default:user;not null"`
	VerifiedAt  *time.Time // nil = ยังไม่ verify
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (UserGormModel) TableName() string { return "users" }
