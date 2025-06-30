package responses

import "github.com/gofiber/fiber/v2"

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"` // Can be CustomError or string
}

type PaginatedAPIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`      // Can be CustomError or string
	Pagination *Pagination `json:"pagination,omitempty"` // Optional pagination info
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// OK returns a successful API response
func OK(c *fiber.Ctx, data interface{}, message string) error {
	if message == "" {
		message = "Operation successful"
	}
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created returns a successful API response for resource creation
func Created(c *fiber.Ctx, data interface{}, message string) error {
	if message == "" {
		message = "Resource created successfully"
	}
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error returns an error API response
func Error(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// PaginatedOK returns a successful API response with pagination
func PaginatedOK(c *fiber.Ctx, data interface{}, pagination *Pagination, message string) error {
	if message == "" {
		message = "Operation successful"
	}
	return c.Status(fiber.StatusOK).JSON(PaginatedAPIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
