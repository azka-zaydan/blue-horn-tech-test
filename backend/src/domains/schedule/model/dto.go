package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// StartVisitRequest defines the request body for starting a visit
type StartVisitRequest struct {
	ID        string  `json:"id" validate:"required,uuid"` // Schedule ID
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// EndVisitRequest defines the request body for ending a visit
type EndVisitRequest struct {
	ID        string  `json:"id" validate:"required,uuid"` // Schedule ID
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// FilterSchedulesRequest defines the request body for filtering schedules
type FilterSchedulesRequest struct {
	Limit  int    `query:"limit" validate:"required,min=1,max=100"`       //
	Page   int    `query:"page" validate:"required,min=1"`                // Page number for pagination
	Offset int    `query:"-"`                                             // Offset for pagination, optional
	Date   string `query:"date" validate:"omitempty,datetime=2006-01-02"` // Date in YYYY-MM-DD format
}

func (r *FilterSchedulesRequest) Validate() error {
	if r.Page < 1 {
		r.Page = 1 // Ensure page is at least 1
	}
	if r.Limit < 1 {
		r.Limit = 10 // Default limit if not set or invalid
	}
	return validator.New().Struct(r)
}

func (r *FilterSchedulesRequest) SetOffset() {
	r.Offset = (r.Page - 1) * r.Limit
}

func (r *FilterSchedulesRequest) String() string {
	return fmt.Sprintf("FilterSchedulesRequest{Limit: %d, Page: %d, Date: %s}", r.Limit, r.Page, r.Date)
}

// PaginatedSchedulesResponse holds schedules with pagination info (simplified, actual Pagination struct moved to responses)
type PaginatedSchedulesResponse struct {
	Data       []Schedule `json:"data"`
	Page       int        `json:"page"`        // Current page number
	PageSize   int        `json:"page_size"`   // Number of items per page
	TotalData  int        `json:"total_data"`  // Total number of items
	TotalPages int        `json:"total_pages"` // Total number of pages
}

func (r *PaginatedSchedulesResponse) String() string {
	return fmt.Sprintf("PaginatedSchedulesResponse{Page: %d, PageSize: %d, Total: %d, TotalPages: %d}", r.Page, r.PageSize, r.TotalData, r.TotalPages)
}

func (r *PaginatedSchedulesResponse) SetTotalPages() {
	if r.TotalData == 0 {
		r.TotalPages = 0 // No pages if no data
		return
	}
	r.TotalPages = (r.TotalData + r.PageSize - 1) / r.PageSize // Ceiling division
}

func (r *StartVisitRequest) Validate() error {
	return validator.New().Struct(r)
}
func (r *EndVisitRequest) Validate() error {
	return validator.New().Struct(r)
}
