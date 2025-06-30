package service_test

import (
	"context"
	mocks "mini-evv-logger-backend/src/domains/task/mocks/repository"
	"mini-evv-logger-backend/src/domains/task/model"
	"mini-evv-logger-backend/src/domains/task/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	mockTaskRepo *mocks.MockTaskRepository
	ctrl         *gomock.Controller
	svc          service.TaskService
)

func initMocks(t *testing.T) {
	ctrl = gomock.NewController(t)

	mockTaskRepo = mocks.NewMockTaskRepository(ctrl)

	svc = service.NewTaskService(mockTaskRepo)
}

func TestGetTasksByScheduleID(t *testing.T) {
	initMocks(t)

	defer ctrl.Finish()

	dummyScheduleID := "test-schedule-id"

	t.Run("TestGetTasksByScheduleID: OK", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTasksByScheduleID(gomock.Any(), dummyScheduleID).Return([]model.Task{
			{ID: "task-id-1", ScheduleID: dummyScheduleID, Description: "Task 1", Status: "pending"},
			{ID: "task-id-2", ScheduleID: dummyScheduleID, Description: "Task 2", Status: "completed"},
		}, nil).Times(1)

		tasks, err := svc.GetTasksBySchedule(context.Background(), dummyScheduleID)
		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
	})

	t.Run("TestGetTasksByScheduleID: No Tasks Found", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTasksByScheduleID(gomock.Any(), dummyScheduleID).Return(nil, nil).Times(1)

		tasks, err := svc.GetTasksBySchedule(context.Background(), dummyScheduleID)
		assert.NoError(t, err)
		assert.Nil(t, tasks)
	})

	t.Run("TestGetTasksByScheduleID: Repository Error", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTasksByScheduleID(gomock.Any(), dummyScheduleID).Return(nil, assert.AnError).Times(1)

		tasks, err := svc.GetTasksBySchedule(context.Background(), dummyScheduleID)
		assert.Error(t, err)
		assert.Nil(t, tasks)
	})
}

func TestUpdateTaskStatus(t *testing.T) {
	initMocks(t)

	defer ctrl.Finish()

	dummyTaskID := uuid.NewString()
	dummyStatus := "completed"
	dummyReason := "Task completed successfully"

	t.Run("TestUpdateTaskStatus: OK", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTaskByID(gomock.Any(), dummyTaskID).Return(&model.Task{ID: dummyTaskID, Status: "pending"}, nil).Times(1)
		mockTaskRepo.EXPECT().UpdateTaskStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := svc.UpdateTaskStatus(context.Background(), model.UpdateTaskStatusRequest{
			TaskID: dummyTaskID,
			Status: dummyStatus,
			Reason: dummyReason,
		})
		assert.NoError(t, err)
	})

	t.Run("TestUpdateTaskStatus: Task Not Found", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTaskByID(gomock.Any(), gomock.Any()).Return(nil, assert.AnError).Times(1)

		err := svc.UpdateTaskStatus(context.Background(), model.UpdateTaskStatusRequest{
			TaskID: dummyTaskID,
			Status: dummyStatus,
			Reason: dummyReason,
		})
		assert.Error(t, err)
	})

	t.Run("TestUpdateTaskStatus: Invalid Status Transition", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTaskByID(gomock.Any(), gomock.Any()).Return(&model.Task{ID: dummyTaskID, Status: "completed"}, nil).Times(1)

		err := svc.UpdateTaskStatus(context.Background(), model.UpdateTaskStatusRequest{
			TaskID: dummyTaskID,
			Status: "pending",
			Reason: dummyReason,
		})
		assert.Error(t, err)
	})
}
