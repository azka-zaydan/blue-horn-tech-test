package controller

import (
	// Import context
	"mini-evv-logger-backend/exceptions"
	"mini-evv-logger-backend/responses"
	"mini-evv-logger-backend/src/domains/task/model" // Import model for DTOs
	"mini-evv-logger-backend/src/domains/task/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// TaskController handles HTTP requests for tasks
type TaskController struct {
	svc service.TaskService
}

// NewTaskController creates a new TaskController
func NewTaskController(svc service.TaskService) *TaskController {
	return &TaskController{svc: svc}
}

// Routes sets up the API endpoints for tasks
func (tc *TaskController) Routes(app fiber.Router) {
	taskRoutes := app.Group("/tasks")
	taskRoutes.Post("/:taskId/update", tc.UpdateTaskStatus)
}

// UpdateTaskStatus handles updating a task's status
func (tc *TaskController) UpdateTaskStatus(c *fiber.Ctx) error {
	ctx := c.UserContext()

	taskID := c.Params("taskId")
	if taskID == "" {
		return responses.Error(c, http.StatusBadRequest, "Task ID is required", exceptions.ErrBadRequest.Error())
	}

	var req model.UpdateTaskStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	req.TaskID = taskID // Set the TaskID from the URL parameter

	err := tc.svc.UpdateTaskStatus(ctx, req)
	if err != nil {
		return exceptions.HandleError(c, err) // Use the new helper
	}
	return responses.OK(c, nil, "Task status updated successfully")
}
