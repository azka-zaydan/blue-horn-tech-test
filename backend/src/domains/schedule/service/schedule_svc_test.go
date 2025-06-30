package service_test

import (
	"context"
	"mini-evv-logger-backend/exceptions"
	mocks "mini-evv-logger-backend/src/domains/schedule/mocks/repository"
	"mini-evv-logger-backend/src/domains/schedule/model"
	"mini-evv-logger-backend/src/domains/schedule/service"

	taskMocks "mini-evv-logger-backend/src/domains/task/mocks/repository"
	taskModel "mini-evv-logger-backend/src/domains/task/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	mockScheduleRepo *mocks.MockScheduleRepository
	mockTaskRepo     *taskMocks.MockTaskRepository
	ctrl             *gomock.Controller
	svc              service.ScheduleService
)

func initMocks(t *testing.T) {
	ctrl = gomock.NewController(t)

	mockScheduleRepo = mocks.NewMockScheduleRepository(ctrl)
	mockTaskRepo = taskMocks.NewMockTaskRepository(ctrl)

	svc = service.NewScheduleService(mockScheduleRepo, mockTaskRepo)
}

func TestGetAllSchedules(t *testing.T) {
	initMocks(t)

	defer ctrl.Finish()

	dummyLimit, dummyPage := 10, 1

	dummyFilter := model.FilterSchedulesRequest{
		Limit: dummyLimit,
		Page:  dummyPage,
		Date:  "",
	}
	t.Run("TestGetAllSchedules: OK", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetSchedules(gomock.Any(), gomock.Any()).Return([]model.Schedule{}, 100, nil).Times(1)

		paginatedSchedules, err := svc.GetAllSchedules(context.Background(), dummyFilter)
		assert.NoError(t, err)
		assert.NotNil(t, paginatedSchedules)
		assert.Equal(t, 100, paginatedSchedules.TotalData)
	})

	t.Run("TestGetAllSchedules: Validation error", func(t *testing.T) {
		invalidFilter := model.FilterSchedulesRequest{
			Limit: 1,
			Page:  dummyPage,
			Date:  "asdasdas",
		}
		_, err := svc.GetAllSchedules(context.Background(), invalidFilter)
		assert.Error(t, err)
	})

	t.Run("TestGetAllSchedules: Error get schedules", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetSchedules(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError).Times(1)
		paginatedSchedules, err := svc.GetAllSchedules(context.Background(), dummyFilter)
		assert.Error(t, err)
		assert.Nil(t, paginatedSchedules)
	})
}

func TestGetScheduleByID(t *testing.T) {
	initMocks(t)

	defer ctrl.Finish()

	dummyID := uuid.NewString()

	t.Run("TestGetScheduleByID: OK", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID}, nil).Times(1)
		mockTaskRepo.EXPECT().GetTasksByScheduleID(gomock.Any(), dummyID).Return([]taskModel.Task{}, nil).Times(1)
		schedule, err := svc.GetScheduleByID(context.Background(), dummyID)
		assert.NoError(t, err)
		assert.NotNil(t, schedule)
		assert.Equal(t, dummyID, schedule.ID)
	})

	t.Run("TestGetScheduleByID: Not Found", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(nil, exceptions.ErrNotFound).Times(1)

		schedule, err := svc.GetScheduleByID(context.Background(), dummyID)
		assert.Error(t, err)
		assert.Nil(t, schedule)
	})

	t.Run("TestGetScheduleByID: Error", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(nil, assert.AnError).Times(1)

		schedule, err := svc.GetScheduleByID(context.Background(), dummyID)
		assert.Error(t, err)
		assert.Nil(t, schedule)
	})

	t.Run("TestGetScheduleByID: Invalid UUID", func(t *testing.T) {
		invalidID := "invalid-uuid"
		schedule, err := svc.GetScheduleByID(context.Background(), invalidID)
		assert.Error(t, err)
		assert.Nil(t, schedule)
		assert.Equal(t, exceptions.ErrBadRequest.WithDetails("Invalid schedule ID format").Error(), err.Error())
	})

	t.Run("TestGetScheduleByID: No Tasks Found", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID}, nil).Times(1)
		mockTaskRepo.EXPECT().GetTasksByScheduleID(gomock.Any(), dummyID).Return(nil, nil).Times(1)
		schedule, err := svc.GetScheduleByID(context.Background(), dummyID)
		assert.NoError(t, err)
		assert.NotNil(t, schedule)
		assert.Equal(t, dummyID, schedule.ID)
		assert.Empty(t, schedule.Tasks) // Ensure no tasks are returned
	})
}

func TestStartVisit(t *testing.T) {
	initMocks(t)

	defer ctrl.Finish()

	dummyID := uuid.NewString()
	dummyRequest := model.StartVisitRequest{
		ID:        dummyID,
		Latitude:  12.345678,
		Longitude: 98.765432,
	}

	t.Run("TestStartVisit: OK", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID, Status: "upcoming"}, nil).Times(1)
		mockScheduleRepo.EXPECT().LogVisitStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := svc.StartVisit(context.Background(), dummyRequest)
		assert.NoError(t, err)
	})

	t.Run("TestStartVisit: Schedule Not Found", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(nil, exceptions.ErrNotFound).Times(1)
		err := svc.StartVisit(context.Background(), dummyRequest)
		assert.Error(t, err)
		assert.Equal(t, exceptions.ErrNotFound.Error(), err.Error())
	})

	t.Run("TestStartVisit: Schedule Already Started", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID, Status: "started"}, nil).Times(1)
		err := svc.StartVisit(context.Background(), dummyRequest)
		assert.Error(t, err)
		assert.Equal(t, exceptions.ErrConflict.WithDetails("Visit for schedule ID "+dummyID+" is already started. Cannot start.").Error(), err.Error())
	})

	t.Run("TestStartVisit: Failed Log Visit Start", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID, Status: "upcoming"}, nil).Times(1)
		mockScheduleRepo.EXPECT().LogVisitStart(gomock.Any(), dummyID, gomock.Any(), dummyRequest.Latitude, dummyRequest.Longitude).Return(assert.AnError).Times(1)

		err := svc.StartVisit(context.Background(), dummyRequest)
		assert.Error(t, err)
	})
}

func TestEndVisit(t *testing.T) {
	initMocks(t)

	defer ctrl.Finish()

	dummyID := uuid.NewString()
	dummyRequest := model.EndVisitRequest{
		ID:        dummyID,
		Latitude:  12.345678,
		Longitude: 98.765432,
	}

	t.Run("TestEndVisit: OK", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID, Status: "in-progress"}, nil).Times(1)
		mockScheduleRepo.EXPECT().LogVisitEnd(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := svc.EndVisit(context.Background(), dummyRequest)
		assert.NoError(t, err)
	})

	t.Run("TestEndVisit: Schedule Not Found", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(nil, exceptions.ErrNotFound).Times(1)
		err := svc.EndVisit(context.Background(), dummyRequest)
		assert.Error(t, err)
		assert.Equal(t, exceptions.ErrNotFound.Error(), err.Error())
	})

	t.Run("TestEndVisit: Schedule Not In Progress", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID, Status: "upcoming"}, nil).Times(1)
		err := svc.EndVisit(context.Background(), dummyRequest)
		assert.Error(t, err)
		assert.Equal(t, exceptions.ErrConflict.WithDetails("Visit for schedule ID "+dummyID+" is currently upcoming. Cannot end.").Error(), err.Error())
	})

	t.Run("TestEndVisit: Failed Log Visit End", func(t *testing.T) {
		mockScheduleRepo.EXPECT().GetScheduleByID(gomock.Any(), dummyID).Return(&model.Schedule{ID: dummyID, Status: "in-progress"}, nil).Times(1)
		mockScheduleRepo.EXPECT().LogVisitEnd(gomock.Any(), dummyID, gomock.Any(), dummyRequest.Latitude, dummyRequest.Longitude).Return(assert.AnError).Times(1)

		err := svc.EndVisit(context.Background(), dummyRequest)
		assert.Error(t, err)
	})

}
