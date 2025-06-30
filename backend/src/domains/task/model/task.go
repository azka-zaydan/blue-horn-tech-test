package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Task represents a care activity/task within a schedule
type Task struct {
	ID          string    `json:"id" db:"id"`
	ScheduleID  string    `json:"schedule_id" db:"schedule_id"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`           // e.g., "pending", "completed", "not_completed"
	Reason      *string   `json:"reason,omitempty" db:"reason"` // Pointer to allow NULL
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateTaskStatusRequest defines the request body for updating a task status
type UpdateTaskStatusRequest struct {
	TaskID string `json:"task_id" validate:"required,uuid"`                                                       // Task ID to update
	Status string `json:"status" validate:"required,oneof=pending completed not_completed cancelled in-progress"` // e.g., "completed", "not_completed"
	Reason string `json:"reason,omitempty" validate:"required_if=status not_completed"`                           // Required if status is "not_completed"
}

func (r *UpdateTaskStatusRequest) Validate() error {
	return validator.New().Struct(r)
}
