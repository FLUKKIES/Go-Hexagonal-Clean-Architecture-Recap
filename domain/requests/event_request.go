package requests

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location" validate:"required"`
	Capacity    int       `json:"capacity" validate:"required,min=1"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}

type JoinEventRequest struct {
	EventID uuid.UUID `json:"event_id" validate:"required"`
}
