package usecases

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/requests"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/responses"
	"github.com/google/uuid"
)

type IEventUsecase interface {
	CreateEvent(adminID uuid.UUID, req *requests.CreateEventRequest) (*responses.EventResponse, error)
	JoinEvent(userID uuid.UUID, req *requests.JoinEventRequest) (*responses.ParticipantResponse, error)
	ListEvents() ([]responses.EventResponse, error)
	GetEventDetails(eventID uuid.UUID) (*responses.EventResponse, error)
}
