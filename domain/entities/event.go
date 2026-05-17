package entities

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title       string    `gorm:"not null"`
	Description string
	Location    string
	Capacity    int       `gorm:"not null"`
	StartTime   time.Time `gorm:"not null"`
	EndTime     time.Time `gorm:"not null"`
	CreatedBy   uuid.UUID `gorm:"not null;type:uuid"` // AdminID
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (Event) TableName() string { return "events" }
