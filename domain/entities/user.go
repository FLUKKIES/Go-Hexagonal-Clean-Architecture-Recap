package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid"`
	FirstName   string     `gorm:"not null"`
	LastName    string     `gorm:"not null"`
	Email       string     `gorm:"uniqueIndex;not null"`
	Password    *string    // nil สำหรับ OAuth users ที่ไม่มี Password
	PhoneNumber *string    `gorm:"uniqueIndex"` // nil จนกว่าจะเพิ่มเบอร์
	ProfileUrl  *string
	Role        UserRole   `gorm:"default:user;not null"` // "user", "admin"
	VerifiedAt  *time.Time // nil = ยังไม่ verify เบอร์มือถือ
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (User) TableName() string { return "users" }
