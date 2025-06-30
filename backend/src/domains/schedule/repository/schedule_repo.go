package repository

import (
	"context" // Import context
	"database/sql"
	"fmt"
	"mini-evv-logger-backend/exceptions"
	"mini-evv-logger-backend/src/domains/schedule/model"
	"time" // Imported for time.Now()

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

//go:generate go run go.uber.org/mock/mockgen -source=./schedule_repo.go -destination=../mocks/repository/schedule_repo.go -package=mocks

// ScheduleRepository defines the interface for schedule database operations
type ScheduleRepository interface {
	GetSchedules(ctx context.Context, filter model.FilterSchedulesRequest) ([]model.Schedule, int, error)
	GetScheduleByID(ctx context.Context, id string) (*model.Schedule, error)
	UpdateScheduleStatus(ctx context.Context, id, status string) error
	LogVisitStart(ctx context.Context, id string, startTime time.Time, latitude, longitude float64) error
	LogVisitEnd(ctx context.Context, id string, endTime time.Time, latitude, longitude float64) error
}

// scheduleRepositoryImpl implements the ScheduleRepository interface
type scheduleRepositoryImpl struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewScheduleRepository creates a new ScheduleRepository (returns interface)
func NewScheduleRepository(db *sqlx.DB, logger zerolog.Logger) ScheduleRepository {
	return &scheduleRepositoryImpl{db: db, logger: logger}
}

// GetSchedules fetches all schedules from the database with pagination
func (r *scheduleRepositoryImpl) GetSchedules(ctx context.Context, filter model.FilterSchedulesRequest) ([]model.Schedule, int, error) {
	var schedules []model.Schedule
	qb := squirrel.Select().
		From("schedules").
		PlaceholderFormat(squirrel.Dollar)

	if filter.Date != "" {
		// If a specific date is provided, filter by shift_time
		qb = qb.Where(
			squirrel.And{
				squirrel.GtOrEq{"shift_time": filter.Date + " 00:00:00"},
				squirrel.LtOrEq{"shift_time": filter.Date + " 23:59:59"},
			},
		)
	}

	countq := qb.Column("COUNT(id)")
	countQuery, countArgs, err := countq.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build SQL query for CountTotalSchedules")
		return nil, 0, exceptions.ErrInternalError
	}

	fmt.Printf("countQuery: %v\n", countQuery)

	var totalCount int
	err = r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Msg("No schedules found in database")
			return []model.Schedule{}, 0, nil // Return empty slice if no schedules exist
		}
		r.logger.Error().Err(err).Msg("Failed to execute SQL query for CountTotalSchedules")
		return nil, 0, exceptions.ErrInternalError
	}

	qb = qb.Columns("id", "client_name", "shift_time", "location", "status",
		"start_time", "start_latitude", "start_longitude", "end_time", "end_latitude", "end_longitude",
		"created_at", "updated_at").
		OrderBy("shift_time ASC").
		Limit(uint64(filter.Limit)).
		Offset(uint64(filter.Offset))

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build SQL query for GetSchedules")
		return nil, 0, exceptions.ErrInternalError
	}

	err = r.db.SelectContext(ctx, &schedules, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Schedule{}, 0, nil
		}
		r.logger.Error().Err(err).Msg("Failed to execute SQL query for GetSchedules")
		return nil, 0, exceptions.ErrInternalError
	}
	return schedules, totalCount, nil
}

// GetScheduleByID fetches a single schedule by ID
func (r *scheduleRepositoryImpl) GetScheduleByID(ctx context.Context, id string) (*model.Schedule, error) {
	var schedule model.Schedule
	qb := squirrel.Select("id", "client_name", "shift_time", "location", "status",
		"start_time", "start_latitude", "start_longitude", "end_time", "end_latitude", "end_longitude",
		"created_at", "updated_at").
		From("schedules").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Msg("Failed to build SQL query for GetScheduleByID")
		return nil, exceptions.ErrInternalError
	}

	err = r.db.GetContext(ctx, &schedule, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("schedule_id", id).Msg("Schedule not found in database")
			return nil, exceptions.ErrNotFound.WithDetails(fmt.Sprintf("Schedule with ID %s not found", id))
		}
		r.logger.Error().Err(err).Str("schedule_id", id).Msg("Failed to execute SQL query for GetScheduleByID")
		return nil, exceptions.ErrInternalError
	}
	return &schedule, nil
}

// UpdateScheduleStatus updates the status of a schedule without pre-checking existence.
// It relies on the service layer to perform existence checks.
func (r *scheduleRepositoryImpl) UpdateScheduleStatus(ctx context.Context, id, status string) error {
	qb := squirrel.Update("schedules").
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Str("status", status).Msg("Failed to build SQL query for UpdateScheduleStatus")
		return exceptions.ErrInternalError
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Str("status", status).Msg("Failed to execute SQL query for UpdateScheduleStatus")
		return exceptions.ErrInternalError
	}

	return nil
}

// LogVisitStart logs the start time and geolocation for a visit.
// It updates the record by ID and sets status to 'in-progress'.
// The service layer is responsible for pre-validating the 'upcoming' status.
func (r *scheduleRepositoryImpl) LogVisitStart(ctx context.Context, id string, startTime time.Time, latitude, longitude float64) error {
	qb := squirrel.Update("schedules").
		Set("start_time", startTime).
		Set("start_latitude", latitude).
		Set("start_longitude", longitude).
		Set("status", "in-progress").
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Msg("Failed to build SQL query for LogVisitStart")
		return exceptions.ErrInternalError
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Msg("Failed to execute SQL query for LogVisitStart")
		return exceptions.ErrInternalError
	}
	return nil
}

// LogVisitEnd logs the end time and geolocation for a visit.
// It updates the record by ID and sets status to 'completed'.
// The service layer is responsible for pre-validating the 'in-progress' status.
func (r *scheduleRepositoryImpl) LogVisitEnd(ctx context.Context, id string, endTime time.Time, latitude, longitude float64) error {
	qb := squirrel.Update("schedules").
		Set("end_time", endTime).
		Set("end_latitude", latitude).
		Set("end_longitude", longitude).
		Set("status", "completed").
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Msg("Failed to build SQL query for LogVisitEnd")
		return exceptions.ErrInternalError
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Str("schedule_id", id).Msg("Failed to execute SQL query for LogVisitEnd")
		return exceptions.ErrInternalError
	}

	return nil
}
