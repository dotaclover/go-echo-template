package common

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse 统一 API 响应结构
type APIResponse struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Data       interface{}            `json:"data,omitempty"`
	Pagination *PaginationMeta        `json:"pagination,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

func requestID(c echo.Context) string {
	if id, ok := c.Get("request_id").(string); ok {
		return id
	}
	return ""
}

func JSON(c echo.Context, httpStatus int, code, message string, data interface{}, pagination *PaginationMeta, details map[string]interface{}) error {
	return c.JSON(httpStatus, APIResponse{
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Details:    details,
		RequestID:  requestID(c),
	})
}

// Success 成功响应
func Success(c echo.Context, message string, data interface{}) error {
	return JSON(c, http.StatusOK, "success", message, data, nil, nil)
}

// Created 创建成功响应
func Created(c echo.Context, message string, data interface{}) error {
	return JSON(c, http.StatusCreated, "created", message, data, nil, nil)
}

// Paginated 分页成功响应
func Paginated(c echo.Context, message string, data interface{}, pagination *PaginationMeta) error {
	return JSON(c, http.StatusOK, "success", message, data, pagination, nil)
}

// Error 错误响应
func Error(c echo.Context, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return JSON(c, appErr.HTTPStatus, appErr.Code, appErr.Message, nil, nil, appErr.Details)
	}
	return JSON(c, http.StatusInternalServerError, "internal_error", "Internal server error", nil, nil, nil)
}

// BadRequest 400 错误响应
func BadRequest(c echo.Context, message string) error {
	return JSON(c, http.StatusBadRequest, "bad_request", message, nil, nil, nil)
}

// Unauthorized 401 错误响应
func Unauthorized(c echo.Context, message string) error {
	return JSON(c, http.StatusUnauthorized, "unauthorized", message, nil, nil, nil)
}

// Forbidden 403 错误响应
func Forbidden(c echo.Context, message string) error {
	return JSON(c, http.StatusForbidden, "forbidden", message, nil, nil, nil)
}

// NotFound 404 错误响应
func NotFound(c echo.Context, message string) error {
	return JSON(c, http.StatusNotFound, "not_found", message, nil, nil, nil)
}
