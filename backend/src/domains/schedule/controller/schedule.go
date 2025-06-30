package controller

import (
	"mini-evv-logger-backend/exceptions"
	"mini-evv-logger-backend/responses"
	"mini-evv-logger-backend/src/domains/schedule/model"
	"mini-evv-logger-backend/src/domains/schedule/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ScheduleController handles HTTP requests for schedules
type ScheduleController struct {
	svc service.ScheduleService
}

// NewScheduleController creates a new ScheduleController
func NewScheduleController(svc service.ScheduleService) *ScheduleController {
	return &ScheduleController{svc: svc}
}

// Routes sets up the API endpoints for schedules
func (sc *ScheduleController) Routes(app fiber.Router) {
	scheduleRoutes := app.Group("/schedules")
	scheduleRoutes.Get("/", sc.GetSchedules)
	scheduleRoutes.Get("/:id", sc.GetScheduleDetails)
	scheduleRoutes.Post("/:id/start", sc.StartVisit)
	scheduleRoutes.Post("/:id/end", sc.EndVisit)
}

// GetSchedules handles fetching all schedules with pagination
func (sc *ScheduleController) GetSchedules(c *fiber.Ctx) error {
	ctx := c.UserContext()

	filter := model.FilterSchedulesRequest{}

	err := c.QueryParser(&filter)
	if err != nil {
		return responses.Error(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	paginatedSchedules, err := sc.svc.GetAllSchedules(ctx, filter)
	if err != nil {
		return exceptions.HandleError(c, err)
	}

	// Construct the Pagination struct for the response
	pagination := &responses.Pagination{
		Page:       paginatedSchedules.Page,
		PageSize:   paginatedSchedules.PageSize,
		TotalItems: paginatedSchedules.TotalData,
		TotalPages: paginatedSchedules.TotalPages,
	}

	return responses.PaginatedOK(c, paginatedSchedules.Data, pagination, "Schedules retrieved successfully")
}

// GetScheduleDetails handles fetching a single schedule's details
func (sc *ScheduleController) GetScheduleDetails(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return responses.Error(c, http.StatusBadRequest, "Schedule ID is required", exceptions.ErrBadRequest.Error())
	}

	schedule, err := sc.svc.GetScheduleByID(ctx, id)
	if err != nil {
		return exceptions.HandleError(c, err)
	}
	return responses.OK(c, schedule, "Schedule details retrieved successfully")
}

// StartVisit handles logging the start of a visit
func (sc *ScheduleController) StartVisit(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return responses.Error(c, http.StatusBadRequest, "Schedule ID is required", exceptions.ErrBadRequest.Error())
	}

	var req model.StartVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	req.ID = id // Set the ID from the URL parameter
	err := sc.svc.StartVisit(ctx, req)
	if err != nil {
		return exceptions.HandleError(c, err)
	}
	return responses.OK(c, nil, "Visit started successfully")
}

// EndVisit handles logging the end of a visit
func (sc *ScheduleController) EndVisit(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return responses.Error(c, http.StatusBadRequest, "Schedule ID is required", exceptions.ErrBadRequest.Error())
	}

	var req model.EndVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	req.ID = id // Set the ID from the URL parameter
	err := sc.svc.EndVisit(ctx, req)
	if err != nil {
		return exceptions.HandleError(c, err)
	}
	return responses.OK(c, nil, "Visit ended successfully")
}
