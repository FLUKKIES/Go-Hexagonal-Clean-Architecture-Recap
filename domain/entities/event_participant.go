package entities

import (
	"time"

	"github.com/google/uuid"
)

type EventParticipant struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	EventID   uuid.UUID `gorm:"not null;type:uuid"`
	UserID    uuid.UUID `gorm:"not null;type:uuid"`
	JoinedAt  time.Time `gorm:"not null"`

	// ความสัมพันธ์ (Associations) เพื่อให้ GORM สร้าง Foreign Key ใน Database ให้อัตโนมัติ
	Event Event `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (EventParticipant) TableName() string { return "event_participants" }
