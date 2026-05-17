package responses

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/google/uuid"
)

type EventResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy    uuid.UUID             `json:"created_by"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	Participants []ParticipantResponse `json:"participants,omitempty"`
}

type ParticipantResponse struct {
	ID       uuid.UUID     `json:"id"`
	EventID  uuid.UUID     `json:"event_id"`
	UserID   uuid.UUID     `json:"user_id"`
	JoinedAt time.Time     `json:"joined_at"`
	User     *UserResponse `json:"user,omitempty"`
}

func ToEventResponse(event *entities.Event) EventResponse {
	return EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		Capacity:    event.Capacity,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		CreatedBy:   event.CreatedBy,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func ToParticipantResponse(participant *entities.EventParticipant) ParticipantResponse {
	resp := ParticipantResponse{
		ID:       participant.ID,
		EventID:  participant.EventID,
		UserID:   participant.UserID,
		JoinedAt: participant.JoinedAt,
	}

	// ถ้ามีการโหลดข้อมูล User มาด้วย (Preload) ให้ map ใส่เข้าไป
	if participant.User.ID != uuid.Nil {
		userResp := ToUserResponse(&participant.User)
		resp.User = &userResp
	}

	return resp
}
