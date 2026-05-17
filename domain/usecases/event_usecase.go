package usecases

import (
	"time"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/exceptions"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/repositories"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/requests"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/responses"
	"github.com/google/uuid"
)

type eventUsecaseImpl struct {
	eventRepo            repositories.IEventRepository
	eventParticipantRepo repositories.IEventParticipantRepository
	userRepo             repositories.IUserRepository
}

func NewEventUsecase(
	eventRepo repositories.IEventRepository,
	eventParticipantRepo repositories.IEventParticipantRepository,
	userRepo repositories.IUserRepository,
) IEventUsecase {
	return &eventUsecaseImpl{
		eventRepo:            eventRepo,
		eventParticipantRepo: eventParticipantRepo,
		userRepo:             userRepo,
	}
}

func (u *eventUsecaseImpl) CreateEvent(adminID uuid.UUID, req *requests.CreateEventRequest) (*responses.EventResponse, error) {
	// 1. ตรวจสอบว่า User เป็น Admin หรือไม่
	user, err := u.userRepo.FindByID(adminID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Role != entities.UserRoleAdmin {
		return nil, exceptions.ErrUnauthorizedAction
	}

	// 2. สร้าง Event Entity
	event := &entities.Event{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Capacity:    req.Capacity,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedBy:   adminID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 3. บันทึกลง Database
	if err := u.eventRepo.Create(event); err != nil {
		return nil, err
	}

	// 4. แปลงกลับเป็น Response
	res := responses.ToEventResponse(event)
	return &res, nil
}

func (u *eventUsecaseImpl) JoinEvent(userID uuid.UUID, req *requests.JoinEventRequest) (*responses.ParticipantResponse, error) {
	// 1. หา Event
	event, err := u.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, exceptions.ErrEventNotFound
	}

	// 2. ตรวจสอบว่าเคยเข้าร่วมหรือยัง
	existing, err := u.eventParticipantRepo.FindByEventAndUser(req.EventID, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, exceptions.ErrAlreadyJoined
	}

	// 3. ตรวจสอบ Capacity
	count, err := u.eventParticipantRepo.CountByEventID(req.EventID)
	if err != nil {
		return nil, err
	}
	if count >= int64(event.Capacity) {
		return nil, exceptions.ErrEventFull
	}

	// 4. สร้าง Participant
	participant := &entities.EventParticipant{
		ID:       uuid.New(),
		EventID:  req.EventID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}

	// 5. บันทึกลง Database
	if err := u.eventParticipantRepo.Create(participant); err != nil {
		return nil, err
	}

	res := responses.ToParticipantResponse(participant)
	return &res, nil
}

func (u *eventUsecaseImpl) ListEvents() ([]responses.EventResponse, error) {
	events, err := u.eventRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []responses.EventResponse
	for _, e := range events {
		res = append(res, responses.ToEventResponse(&e))
	}
	return res, nil
}

func (u *eventUsecaseImpl) GetEventDetails(eventID uuid.UUID) (*responses.EventResponse, error) {
	event, err := u.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, exceptions.ErrEventNotFound
	}

	res := responses.ToEventResponse(event)

	// Fetch participants
	participants, err := u.eventParticipantRepo.FindByEventID(eventID)
	if err == nil && len(participants) > 0 {
		for _, p := range participants {
			pResp := responses.ToParticipantResponse(&p)
			res.Participants = append(res.Participants, pResp)
		}
	}

	return &res, nil
}
