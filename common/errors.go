package common

import "net/http"

// AppError 应用层统一错误
// 通过 HTTPStatus 和 Code 将业务错误映射为标准 API 响应。
type AppError struct {
	HTTPStatus int                    `json:"-"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(httpStatus int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

func NewAppErrorWithDetails(httpStatus int, code, message string, details map[string]interface{}) *AppError {
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

func BadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, "bad_request", message)
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, "unauthorized", message)
}

func ForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, "forbidden", message)
}

func NotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, "not_found", message)
}

func ConflictError(message string) *AppError {
	return NewAppError(http.StatusConflict, "conflict", message)
}

func ValidationError(message string, details map[string]interface{}) *AppError {
	return NewAppErrorWithDetails(http.StatusBadRequest, "validation_error", message, details)
}

func InternalError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, "internal_error", message)
}

func TooManyRequestsError(message string) *AppError {
	return NewAppError(http.StatusTooManyRequests, "too_many_requests", message)
}
