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
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Password    *string // nil สำหรับ OAuth users ที่ไม่มี Password
	PhoneNumber *string // nil จนกว่าจะเพิ่มเบอร์
	ProfileUrl  *string
	Role        UserRole   // "user", "admin"
	VerifiedAt  *time.Time // nil = ยังไม่ verify เบอร์มือถือ
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
