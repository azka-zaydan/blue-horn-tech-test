package repository_test

import (
	"context"
	"database/sql"
	pkgmock "mini-evv-logger-backend/pkg_mock"
	"mini-evv-logger-backend/src/domains/schedule/model"
	"mini-evv-logger-backend/src/domains/schedule/repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	dbMock   *sql.DB
	sqlxMock *sqlx.DB
	mockSQL  sqlmock.Sqlmock
	repo     repository.ScheduleRepository
)

func initMocks(t *testing.T) {
	var err error
	dbMock, mockSQL, err = sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	// Wrap sqlmock in sqlx.DB
	sqlxMock = sqlx.NewDb(dbMock, "sqlmock")
	repo = repository.NewScheduleRepository(sqlxMock, pkgmock.InitMockLogger())
}

func TestGetSchedules(t *testing.T) {
	initMocks(t)

	dummyLimit, dummyOffset := 10, 0

	countQuery := `SELECT COUNT(id) FROM schedules`
	query := `SELECT id, client_name, shift_time, location, status, start_time, start_latitude, start_longitude, end_time, end_latitude, end_longitude, created_at, updated_at FROM schedules ORDER BY shift_time ASC LIMIT 10 OFFSET 0`
	dummySchedules := []model.Schedule{
		{
			ID:             uuid.NewString(),
			ClientName:     "Test Client",
			ShiftTime:      time.Now(),
			Location:       "Test Location",
			Status:         "scheduled",
			StartTime:      nil,
			StartLatitude:  nil,
			StartLongitude: nil,
			EndTime:        nil,
			EndLatitude:    nil,
			EndLongitude:   nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	dummyFilter := model.FilterSchedulesRequest{
		Limit:  dummyLimit,
		Offset: dummyOffset,
		Date:   "",
	}

	t.Run("TestGetSchedules: OK", func(t *testing.T) {

		mockSQL.ExpectQuery(regexp.QuoteMeta(countQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(dummySchedules)))
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "client_name", "shift_time", "location", "status", "start_time", "start_latitude", "start_longitude", "end_time", "end_latitude", "end_longitude", "created_at", "updated_at"}).
				AddRow(dummySchedules[0].ID, dummySchedules[0].ClientName, dummySchedules[0].ShiftTime, dummySchedules[0].Location, dummySchedules[0].Status, dummySchedules[0].StartTime, dummySchedules[0].StartLatitude, dummySchedules[0].StartLongitude, dummySchedules[0].EndTime, dummySchedules[0].EndLatitude, dummySchedules[0].EndLongitude, dummySchedules[0].CreatedAt, dummySchedules[0].UpdatedAt))

		schedules, total, err := repo.GetSchedules(context.Background(), dummyFilter)
		assert.Nil(t, err)
		assert.Equal(t, len(dummySchedules), total)
		assert.Len(t, schedules, 1)
	})

	t.Run("TestGetSchedules: No Rows", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(countQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(dummySchedules)))
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrNoRows)
		schedules, total, err := repo.GetSchedules(context.Background(), dummyFilter)
		assert.Nil(t, err)
		assert.Equal(t, 0, total)
		assert.Len(t, schedules, 0)
	})

	t.Run("TestGetSchedules: SQL Error", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)
		schedules, total, err := repo.GetSchedules(context.Background(), dummyFilter)
		assert.NotNil(t, err)
		assert.Equal(t, "Error 500: Internal server error", err.Error())
		assert.Len(t, schedules, 0)
		assert.Equal(t, 0, total)
	})

}

func TestGetScheduleByID(t *testing.T) {
	initMocks(t)

	dummyID := uuid.NewString()
	query := `SELECT id, client_name, shift_time, location, status, start_time, start_latitude, start_longitude, end_time, end_latitude, end_longitude, created_at, updated_at FROM schedules WHERE id = $1`
	dummySchedule := model.Schedule{
		ID:             dummyID,
		ClientName:     "Test Client",
		ShiftTime:      time.Now(),
		Location:       "Test Location",
		Status:         "scheduled",
		StartTime:      nil,
		StartLatitude:  nil,
		StartLongitude: nil,
		EndTime:        nil,
		EndLatitude:    nil,
		EndLongitude:   nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	t.Run("TestGetScheduleByID: OK", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(dummyID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "client_name", "shift_time", "location", "status", "start_time", "start_latitude", "start_longitude", "end_time", "end_latitude", "end_longitude", "created_at", "updated_at"}).
				AddRow(dummySchedule.ID, dummySchedule.ClientName, dummySchedule.ShiftTime, dummySchedule.Location, dummySchedule.Status, dummySchedule.StartTime, dummySchedule.StartLatitude, dummySchedule.StartLongitude, dummySchedule.EndTime, dummySchedule.EndLatitude, dummySchedule.EndLongitude, dummySchedule.CreatedAt, dummySchedule.UpdatedAt))

		schedule, err := repo.GetScheduleByID(context.Background(), dummyID)
		assert.Nil(t, err)
		assert.Equal(t, &dummySchedule, schedule)
	})

	t.Run("TestGetScheduleByID: No Rows", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(dummyID).
			WillReturnError(sql.ErrNoRows)

		schedule, err := repo.GetScheduleByID(context.Background(), dummyID)
		assert.NotNil(t, err)
		assert.Nil(t, schedule)
	})

	t.Run("TestGetScheduleByID: SQL Error", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(dummyID).
			WillReturnError(sql.ErrConnDone)

		schedule, err := repo.GetScheduleByID(context.Background(), dummyID)
		assert.NotNil(t, err)
		assert.Equal(t, "Error 500: Internal server error", err.Error())
		assert.Nil(t, schedule)
	})
}

func TestUpdateScheduleStatus(t *testing.T) {
	initMocks(t)

	dummyID := uuid.NewString()
	dummyStatus := "completed"
	query := `UPDATE schedules SET status = $1, updated_at = $2 WHERE id
    = $3`
	t.Run("TestUpdateScheduleStatus: OK", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(dummyStatus, sqlmock.AnyArg(), dummyID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateScheduleStatus(context.Background(), dummyID, dummyStatus)
		assert.Nil(t, err)
	})

	t.Run("TestUpdateScheduleStatus: SQL Error", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateScheduleStatus(context.Background(), dummyID, dummyStatus)
		assert.NotNil(t, err)
		assert.Equal(t, "Error 500: Internal server error", err.Error())
	})

}

func TestLogVisitStart(t *testing.T) {
	initMocks(t)

	dummyID := uuid.NewString()
	dummyStartTime := time.Now()
	dummyLatitude, dummyLongitude := 37.7749, -122.4194
	query := `UPDATE schedules SET start_time = $1, start_latitude = $2, start_longitude = $3, status = $4, updated_at = $5 WHERE id = $6`
	t.Run("TestLogVisitStart: OK", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(dummyStartTime, dummyLatitude, dummyLongitude, sqlmock.AnyArg(), sqlmock.AnyArg(), dummyID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogVisitStart(context.Background(), dummyID, dummyStartTime, dummyLatitude, dummyLongitude)
		assert.Nil(t, err)
	})

	t.Run("TestLogVisitStart: SQL Error", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)
		err := repo.LogVisitStart(context.Background(), dummyID, dummyStartTime, dummyLatitude, dummyLongitude)
		assert.NotNil(t, err)
		assert.Equal(t, "Error 500: Internal server error", err.Error())
	})
}

func TestLogVisitEnd(t *testing.T) {
	initMocks(t)

	dummyID := uuid.NewString()
	dummyEndTime := time.Now()
	dummyLatitude, dummyLongitude := 37.7749, -122.4194
	query := `UPDATE schedules SET end_time = $1, end_latitude = $2, end_longitude = $3, status = $4, updated_at = $5 WHERE id = $6`
	t.Run("TestLogVisitEnd: OK", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(dummyEndTime, dummyLatitude, dummyLongitude, "completed", sqlmock.AnyArg(), dummyID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogVisitEnd(context.Background(), dummyID, dummyEndTime, dummyLatitude, dummyLongitude)
		assert.Nil(t, err)
	})

	t.Run("TestLogVisitEnd: SQL Error", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)
		err := repo.LogVisitEnd(context.Background(), dummyID, dummyEndTime, dummyLatitude, dummyLongitude)
		assert.NotNil(t, err)
		assert.Equal(t, "Error 500: Internal server error", err.Error())
	})
}
