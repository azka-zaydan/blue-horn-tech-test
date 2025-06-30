package service

import (
	"context" // Import context
	"fmt"
	"mini-evv-logger-backend/exceptions"
	"mini-evv-logger-backend/src/domains/schedule/model"
	"mini-evv-logger-backend/src/domains/schedule/repository"
	taskRepo "mini-evv-logger-backend/src/domains/task/repository"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// ScheduleService defines the interface for schedule business logic
type ScheduleService interface {
	GetAllSchedules(ctx context.Context, filter model.FilterSchedulesRequest) (*model.PaginatedSchedulesResponse, error)
	GetScheduleByID(ctx context.Context, id string) (*model.Schedule, error)
	StartVisit(ctx context.Context, req model.StartVisitRequest) error
	EndVisit(ctx context.Context, req model.EndVisitRequest) error
}

// scheduleServiceImpl implements the ScheduleService interface
type scheduleServiceImpl struct {
	scheduleRepo repository.ScheduleRepository
	taskRepo     taskRepo.TaskRepository
}

// NewScheduleService creates a new ScheduleService (returns interface)
func NewScheduleService(scheduleRepo repository.ScheduleRepository, taskRepo taskRepo.TaskRepository) ScheduleService {
	return &scheduleServiceImpl{scheduleRepo: scheduleRepo, taskRepo: taskRepo}
}

// GetAllSchedules fetches all schedules with pagination
func (s *scheduleServiceImpl) GetAllSchedules(ctx context.Context, filter model.FilterSchedulesRequest) (*model.PaginatedSchedulesResponse, error) {
	log.Info().Msgf("Fetching all schedules with filter %s", filter.String())

	err := filter.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for FilterSchedulesRequest")
		return nil, exceptions.ErrBadRequest.WithDetails(err.Error())
	}

	filter.SetOffset()

	schedules, total, err := s.scheduleRepo.GetSchedules(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all schedules from repository")
		return nil, err
	}

	res := model.PaginatedSchedulesResponse{
		Data:      schedules,
		TotalData: total,
		Page:      filter.Page,
		PageSize:  filter.Limit,
	}
	res.SetTotalPages()

	return &res, nil
}

// GetScheduleByID fetches a schedule by its ID, including its associated tasks
func (s *scheduleServiceImpl) GetScheduleByID(ctx context.Context, id string) (*model.Schedule, error) {
	log.Info().Str("schedule_id", id).Msg("Fetching schedule by ID")

	_, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", id).Msg("Invalid UUID format for schedule ID")
		return nil, exceptions.ErrBadRequest.WithDetails("Invalid schedule ID format")
	}
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", id).Msg("Failed to fetch schedule by ID from repository")
		return nil, err
	}

	tasks, err := s.taskRepo.GetTasksByScheduleID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", id).Msg("Failed to fetch tasks for schedule")
		return schedule, nil
	}
	schedule.Tasks = tasks
	return schedule, nil
}

// StartVisit updates the schedule with start time and geolocation
func (s *scheduleServiceImpl) StartVisit(ctx context.Context, req model.StartVisitRequest) error {
	log.Info().Str("schedule_id", req.ID).Float64("latitude", req.Latitude).Float64("longitude", req.Longitude).Msg("Attempting to start visit")

	err := req.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for StartVisitRequest")
		return exceptions.ErrBadRequest.WithDetails(err.Error())
	}

	// 1. Check if the schedule exists and its current status
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, req.ID)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", req.ID).Msg("Failed to retrieve schedule before starting visit")
		return err
	}

	// 2. Apply business logic: only "upcoming" schedules can be started
	if schedule.Status != "upcoming" {
		return exceptions.ErrConflict.WithDetails(fmt.Sprintf("Visit for schedule ID %s is already %s. Cannot start.", req.ID, schedule.Status))
	}

	// 3. Perform the update via repository
	err = s.scheduleRepo.LogVisitStart(ctx, req.ID, time.Now(), req.Latitude, req.Longitude)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", req.ID).Msg("Failed to log visit start in repository")
		return err
	}
	return nil
}

// EndVisit updates the schedule with end time and geolocation
func (s *scheduleServiceImpl) EndVisit(ctx context.Context, req model.EndVisitRequest) error {
	log.Info().Str("schedule_id", req.ID).Float64("latitude", req.Latitude).Float64("longitude", req.Longitude).Msg("Attempting to end visit")

	err := req.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for EndVisitRequest")
		return exceptions.ErrBadRequest.WithDetails(err.Error())
	}

	// 1. Check if the schedule exists and its current status
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, req.ID)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", req.ID).Msg("Failed to retrieve schedule before ending visit")
		return err
	}

	// 2. Apply business logic: only "in-progress" schedules can be ended
	if schedule.Status != "in-progress" {
		return exceptions.ErrConflict.WithDetails(fmt.Sprintf("Visit for schedule ID %s is currently %s. Cannot end.", req.ID, schedule.Status))
	}

	// 3. Perform the update via repository
	err = s.scheduleRepo.LogVisitEnd(ctx, req.ID, time.Now(), req.Latitude, req.Longitude)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", req.ID).Msg("Failed to log visit end in repository")
		return err
	}
	return nil
}

// UpdateScheduleStatus handles updating the status of a schedule.
func (s *scheduleServiceImpl) UpdateScheduleStatus(ctx context.Context, id, status string) error {
	log.Info().Str("schedule_id", id).Str("status", status).Msg("Attempting to update schedule status")

	// 1. Check if the schedule exists
	_, err := s.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", id).Msg("Schedule not found for status update")
		return err
	}

	// 2. Perform the update
	err = s.scheduleRepo.UpdateScheduleStatus(ctx, id, status)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", id).Str("status", status).Msg("Failed to update schedule status in repository")
		return err
	}
	return nil
}
