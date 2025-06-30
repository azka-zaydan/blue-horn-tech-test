package model

import (
	taskModel "mini-evv-logger-backend/src/domains/task/model"
	"time"
)

// Schedule represents a caregiver's schedule
type Schedule struct {
	ID             string           `json:"id" db:"id"`
	ClientName     string           `json:"client_name" db:"client_name"`
	ShiftTime      time.Time        `json:"shift_time" db:"shift_time"`
	Location       string           `json:"location" db:"location"`
	Status         string           `json:"status" db:"status"`                   // e.g., "upcoming", "in-progress", "completed", "missed"
	StartTime      *time.Time       `json:"start_time" db:"start_time"`           // Pointer to allow NULL
	StartLatitude  *float64         `json:"start_latitude" db:"start_latitude"`   // Pointer to allow NULL
	StartLongitude *float64         `json:"start_longitude" db:"start_longitude"` // Pointer to allow NULL
	EndTime        *time.Time       `json:"end_time" db:"end_time"`               // Pointer to allow NULL
	EndLatitude    *float64         `json:"end_latitude" db:"end_latitude"`       // Pointer to allow NULL
	EndLongitude   *float64         `json:"end_longitude" db:"end_longitude"`     // Pointer to allow NULL
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
	Tasks          []taskModel.Task `json:"tasks,omitempty" db:"-"` // For schedule details, includes associated tasks
}
