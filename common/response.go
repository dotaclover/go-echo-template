package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse 统一 API 响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedData 分页数据
type PaginatedData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

// Success 成功响应
func Success(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Code:    201,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400 错误响应
func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Code:    400,
		Message: message,
	})
}

// Unauthorized 401 错误响应
func Unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Code:    401,
		Message: message,
	})
}

// Forbidden 403 错误响应
func Forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Code:    403,
		Message: message,
	})
}

// NotFound 404 错误响应
func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Code:    404,
		Message: message,
	})
}

// InternalError 500 错误响应
func InternalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Code:    500,
		Message: message,
	})
}
