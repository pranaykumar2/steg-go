package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is a standardized API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse returns a standardized success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse returns a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// ValidationErrorResponse returns a standardized validation error response
func ValidationErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// InternalErrorResponse returns a standardized internal server error response
func InternalErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal server error: "+err.Error())
}

// NotFoundResponse returns a standardized not found error response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// UnauthorizedResponse returns a standardized unauthorized error response
func UnauthorizedResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusUnauthorized, "Unauthorized access")
}
