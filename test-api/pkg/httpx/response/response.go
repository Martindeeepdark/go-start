package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess            = 0
	CodeInvalidParams      = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeInternalError      = 500
	CodeServiceUnavailable = 503
)

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedData represents paginated response data
type PaginatedData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// Success sends a successful response with data
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage sends a successful response with custom message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response with code and message
func Error(c *gin.Context, code int, message string) {
	statusCode := http.StatusOK
	if code >= CodeInternalError {
		statusCode = http.StatusInternalServerError
	} else if code == CodeNotFound {
		statusCode = http.StatusNotFound
	} else if code == CodeUnauthorized {
		statusCode = http.StatusUnauthorized
	} else if code == CodeForbidden {
		statusCode = http.StatusForbidden
	} else if code >= CodeInvalidParams && code < CodeUnauthorized {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData sends an error response with data
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	statusCode := http.StatusOK
	if code >= CodeInternalError {
		statusCode = http.StatusInternalServerError
	} else if code == CodeNotFound {
		statusCode = http.StatusNotFound
	} else if code == CodeUnauthorized {
		statusCode = http.StatusUnauthorized
	} else if code == CodeForbidden {
		statusCode = http.StatusForbidden
	} else if code >= CodeInvalidParams && code < CodeUnauthorized {
		statusCode = http.StatusBadRequest
	}

	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request error
func BadRequest(c *gin.Context, message string) {
	Error(c, CodeInvalidParams, message)
}

// Unauthorized sends an unauthorized error
func Unauthorized(c *gin.Context, message string) {
	Error(c, CodeUnauthorized, message)
}

// Forbidden sends a forbidden error
func Forbidden(c *gin.Context, message string) {
	Error(c, CodeForbidden, message)
}

// NotFound sends a not found error
func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

// InternalError sends an internal server error
func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalError, message)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	Success(c, PaginatedData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
