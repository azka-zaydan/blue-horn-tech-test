package repository_test

import (
	"context"
	"database/sql"
	pkgmock "mini-evv-logger-backend/pkg_mock"
	"mini-evv-logger-backend/src/domains/task/repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	dbMock   *sql.DB
	sqlxMock *sqlx.DB
	mockSQL  sqlmock.Sqlmock
	repo     repository.TaskRepository
)

func initMocks(t *testing.T) {
	var err error
	dbMock, mockSQL, err = sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	// Wrap sqlmock in sqlx.DB
	sqlxMock = sqlx.NewDb(dbMock, "sqlmock")
	repo = repository.NewTaskRepository(sqlxMock, pkgmock.InitMockLogger())
}

func TestGetTasksByScheduleID(t *testing.T) {
	initMocks(t)

	// Define the expected query and result
	scheduleID := "test-schedule-id"

	query := "SELECT id, schedule_id, description, status, reason, created_at, updated_at FROM tasks WHERE schedule_id = $1 ORDER BY created_at ASC"

	t.Run("TestGetTasksByScheduleID: OK", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(scheduleID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "schedule_id", "description", "status", "reason", "created_at", "updated_at"}).
				AddRow("task-id-1", scheduleID, "Task 1", "pending", nil, time.Now(), time.Now()).
				AddRow("task-id-2", scheduleID, "Task 2", "completed", nil, time.Now(), time.Now()))
		tasks, err := repo.GetTasksByScheduleID(context.Background(), scheduleID)
		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
	})

	t.Run("TestGetTasksByScheduleID: No Rows", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(scheduleID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "schedule_id", "description", "status", "reason", "created_at", "updated_at"}))
		tasks, err := repo.GetTasksByScheduleID(context.Background(), scheduleID)
		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
	})

	t.Run("TestGetTasksByScheduleID: Error", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(scheduleID).
			WillReturnError(sql.ErrConnDone)
		tasks, err := repo.GetTasksByScheduleID(context.Background(), scheduleID)
		assert.Error(t, err)
		assert.Nil(t, tasks)
	})

}

func TestGetTaskByID(t *testing.T) {
	initMocks(t)

	// Define the expected query and result
	taskID := "test-task-id"

	query := "SELECT id, schedule_id, description, status, reason, created_at, updated_at FROM tasks WHERE id = $1"

	t.Run("TestGetTaskByID: OK", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(taskID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "schedule_id", "description", "status", "reason", "created_at", "updated_at"}).
				AddRow(taskID, "test-schedule-id", "Test Task", "pending", nil, time.Now(), time.Now()))
		task, err := repo.GetTaskByID(context.Background(), taskID)
		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, taskID, task.ID)
	})

	t.Run("TestGetTaskByID: No Rows", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(taskID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "schedule_id", "description", "status", "reason", "created_at", "updated_at"}))
		task, err := repo.GetTaskByID(context.Background(), taskID)
		assert.Error(t, err)
		assert.Nil(t, task)
	})

	t.Run("TestGetTaskByID: Error", func(t *testing.T) {
		mockSQL.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(taskID).
			WillReturnError(sql.ErrConnDone)
		task, err := repo.GetTaskByID(context.Background(), taskID)
		assert.Error(t, err)
		assert.Nil(t, task)
	})
}

func TestUpdateTaskStatus(t *testing.T) {
	initMocks(t)

	// Define the expected query and result
	taskID := "test-task-id"
	status := "completed"
	reason := "Task completed successfully"

	query := "UPDATE tasks SET status = $1, reason = $2, updated_at = $3 WHERE id = $4"

	t.Run("TestUpdateTaskStatus: OK", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(status, reason, sqlmock.AnyArg(), taskID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.UpdateTaskStatus(context.Background(), taskID, status, &reason)
		assert.NoError(t, err)
	})

	t.Run("TestUpdateTaskStatus: Error", func(t *testing.T) {
		mockSQL.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(status, reason, sqlmock.AnyArg(), taskID).
			WillReturnError(sql.ErrConnDone)
		err := repo.UpdateTaskStatus(context.Background(), taskID, status, &reason)
		assert.Error(t, err)
	})
}
