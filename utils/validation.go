package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ValidationService 验证服务
type ValidationService struct {
	validator *validator.Validate
}

// NewValidationService 创建验证服务实例
func NewValidationService() *ValidationService {
	return &ValidationService{validator: validator.New()}
}

// ValidateStruct 验证结构体
func (v *ValidationService) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

// ValidateIDParam 验证路由中的ID参数
func (v *ValidationService) ValidateIDParam(c echo.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, errors.New("missing ID parameter")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		return 0, errors.New("invalid ID")
	}
	return uint(id), nil
}

// ValidatePagination 验证分页参数
func (v *ValidationService) ValidatePagination(page, limit string) (int, int, error) {
	pageNum := 1
	limitNum := 10

	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil || p < 1 {
			return 0, 0, errors.New("invalid page")
		}
		pageNum = p
	}
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil || l < 1 || l > 100 {
			return 0, 0, errors.New("invalid limit (1-100)")
		}
		limitNum = l
	}
	return pageNum, limitNum, nil
}

// GetErrorMessage 获取验证错误的可读消息
func GetErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// FormatValidationErrors 格式化验证错误为字符串
func FormatValidationErrors(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var msgs []string
		for _, fe := range validationErrs {
			msgs = append(msgs, GetErrorMessage(fe))
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}

// Validator 全局验证服务实例
var Validator = NewValidationService()
