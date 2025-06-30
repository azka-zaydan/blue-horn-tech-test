package repository

import (
	"context" // Import context
	"database/sql"
	"fmt"
	"mini-evv-logger-backend/exceptions"
	"mini-evv-logger-backend/src/domains/task/model"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

//go:generate go run go.uber.org/mock/mockgen -source=./task_repo.go -destination=../mocks/repository/task_repo.go -package=mocks

// TaskRepository defines the interface for task database operations
type TaskRepository interface {
	GetTasksByScheduleID(ctx context.Context, scheduleID string) ([]model.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*model.Task, error)
	UpdateTaskStatus(ctx context.Context, taskID, status string, reason *string) error
}

// taskRepositoryImpl implements the TaskRepository interface
type taskRepositoryImpl struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewTaskRepository creates a new TaskRepository (returns interface)
func NewTaskRepository(db *sqlx.DB, logger zerolog.Logger) TaskRepository {
	return &taskRepositoryImpl{db: db, logger: logger}
}

// GetTasksByScheduleID fetches tasks for a given schedule
func (r *taskRepositoryImpl) GetTasksByScheduleID(ctx context.Context, scheduleID string) ([]model.Task, error) {
	var tasks []model.Task
	qb := squirrel.Select("id", "schedule_id", "description", "status", "reason",
		"created_at", "updated_at").
		From("tasks").
		Where(squirrel.Eq{"schedule_id": scheduleID}).
		OrderBy("created_at ASC").
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", scheduleID).Msg("Failed to build SQL query for GetTasksByScheduleID")
		return nil, exceptions.ErrInternalError
	}

	// Use SelectContext
	err = r.db.SelectContext(ctx, &tasks, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("schedule_id", scheduleID).Msg("No tasks found for this schedule")
			return []model.Task{}, nil
		}
		r.logger.Error().Err(err).Str("schedule_id", scheduleID).Msg("Failed to execute SQL query for GetTasksByScheduleID")
		return nil, exceptions.ErrInternalError
	}
	return tasks, nil
}

// GetTaskByID fetches a single task by ID
func (r *taskRepositoryImpl) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	var task model.Task
	qb := squirrel.Select("id", "schedule_id", "description", "status", "reason",
		"created_at", "updated_at").
		From("tasks").
		Where(squirrel.Eq{"id": taskID}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to build SQL query for GetTaskByID")
		return nil, exceptions.ErrInternalError
	}

	// Use GetContext
	err = r.db.GetContext(ctx, &task, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("task_id", taskID).Msg("Task not found in database")
			return nil, exceptions.ErrNotFound.WithDetails(fmt.Sprintf("Task with ID %s not found", taskID))
		}
		r.logger.Error().Err(err).Str("task_id", taskID).Msg("Failed to execute SQL query for GetTaskByID")
		return nil, exceptions.ErrInternalError
	}
	return &task, nil
}

// UpdateTaskStatus updates the status and optional reason for a task without pre-checking existence.
// It relies on the service layer to perform existence checks.
func (r *taskRepositoryImpl) UpdateTaskStatus(ctx context.Context, taskID, status string, reason *string) error {
	qb := squirrel.Update("tasks").
		Set("status", status).
		Set("reason", reason).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": taskID}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("task_id", taskID).Str("status", status).Msg("Failed to build SQL query for UpdateTaskStatus")
		return exceptions.ErrInternalError
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Str("task_id", taskID).Str("status", status).Msg("Failed to execute SQL query for UpdateTaskStatus")
		return exceptions.ErrInternalError
	}

	return nil
}
