package service

import (
	"context" // Import context
	"fmt"
	"mini-evv-logger-backend/exceptions"
	"mini-evv-logger-backend/src/domains/task/model"
	"mini-evv-logger-backend/src/domains/task/repository"

	"github.com/rs/zerolog/log"
)

// TaskService defines the interface for task business logic
type TaskService interface {
	GetTasksBySchedule(ctx context.Context, scheduleID string) ([]model.Task, error)
	UpdateTaskStatus(ctx context.Context, req model.UpdateTaskStatusRequest) error
}

// taskServiceImpl implements the TaskService interface
type taskServiceImpl struct {
	repo repository.TaskRepository
}

// NewTaskService creates a new TaskService (returns interface)
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskServiceImpl{repo: repo}
}

// GetTasksBySchedule fetches tasks for a specific schedule
func (s *taskServiceImpl) GetTasksBySchedule(ctx context.Context, scheduleID string) ([]model.Task, error) {
	log.Info().Str("schedule_id", scheduleID).Msg("Fetching tasks by schedule ID")
	tasks, err := s.repo.GetTasksByScheduleID(ctx, scheduleID)
	if err != nil {
		log.Error().Err(err).Str("schedule_id", scheduleID).Msg("Failed to fetch tasks by schedule ID from repository")
		return nil, err
	}
	return tasks, nil
}

// UpdateTaskStatus updates a task's status and optional reason
func (s *taskServiceImpl) UpdateTaskStatus(ctx context.Context, req model.UpdateTaskStatusRequest) error {
	log.Info().Str("task_id", req.TaskID).Str("status", req.Status).Str("reason", req.Reason).Msg("Attempting to update task status")

	err := req.Validate()
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for UpdateTaskStatusRequest")
		return exceptions.ErrBadRequest.WithDetails(err.Error())
	}

	// 1. Check if the task exists
	task, err := s.repo.GetTaskByID(ctx, req.TaskID)
	if err != nil {
		log.Error().Err(err).Str("task_id", req.TaskID).Msg("Failed to retrieve task before updating status")
		return err
	}

	// 2. Apply business logic: validate status transition if needed
	if task.Status == "completed" && req.Status == "pending" {
		return exceptions.ErrConflict.WithDetails(fmt.Sprintf("Task ID %s is already completed. Cannot change to pending.", req.TaskID))
	}

	var reasonPtr *string
	if req.Reason != "" {
		reasonPtr = &req.Reason
	}

	// 3. Perform the update via repository
	err = s.repo.UpdateTaskStatus(ctx, req.TaskID, req.Status, reasonPtr)
	if err != nil {
		log.Error().Err(err).Str("task_id", req.TaskID).Msg("Failed to update task status in repository")
		return err
	}
	return nil
}
