package rest

import (
	"errors"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/exceptions"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/requests"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/usecases"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventController struct {
	usecase usecases.IEventUsecase
}

func NewEventController(uc usecases.IEventUsecase) *EventController {
	return &EventController{usecase: uc}
}

// POST /api/events
func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var req requests.CreateEventRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if errs := validator.Validate(&req); errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": errs})
	}

	adminIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	resp, err := c.usecase.CreateEvent(adminID, &req)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// POST /api/events/:id/join
func (c *EventController) JoinEvent(ctx *fiber.Ctx) error {
	eventIDStr := ctx.Params("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid event id"})
	}

	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
	}

	req := requests.JoinEventRequest{EventID: eventID}

	resp, err := c.usecase.JoinEvent(userID, &req)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// GET /api/events
func (c *EventController) ListEvents(ctx *fiber.Ctx) error {
	resp, err := c.usecase.ListEvents()
	if err != nil {
		return c.handleError(ctx, err)
	}

	return ctx.JSON(resp)
}

// GET /api/events/:id
func (c *EventController) GetEventDetails(ctx *fiber.Ctx) error {
	eventIDStr := ctx.Params("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid event id"})
	}

	resp, err := c.usecase.GetEventDetails(eventID)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return ctx.JSON(resp)
}

func (c *EventController) handleError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, exceptions.ErrEventNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrEventFull):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrAlreadyJoined):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrUnauthorizedAction):
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, exceptions.ErrUserNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
