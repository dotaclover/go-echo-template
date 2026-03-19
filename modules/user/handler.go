package user

import (
	"fmt"
	"myapp/models"
	"myapp/server/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler 用户处理器
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes 注册用户相关路由
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	auth := e.Group("/api/auth")

	// 公开路由
	auth.POST("/login", h.Login)

	// 需要认证的路由
	protected := auth.Group("")
	protected.Use(middlewares.JWTAuth())
	protected.GET("/profile", h.GetProfile)
	protected.PUT("/profile", h.UpdateProfile)
	protected.POST("/password", h.UpdatePassword)

	// 管理员路由
	admin := e.Group("/api/admin/users")
	admin.Use(middlewares.JWTAuth(), middlewares.AdminOnly())
	admin.GET("", h.ListUsers)
	admin.POST("", h.CreateUser)
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	token, user, err := h.service.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  toResponse(user),
	})
}

func (h *Handler) GetProfile(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Profile updated",
		"user":    toResponse(user),
	})
}

func (h *Handler) UpdatePassword(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := h.service.UpdatePassword(userID, &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password updated"})
}

func (h *Handler) ListUsers(c echo.Context) error {
	page, pageSize := 1, 20
	if p := c.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	users, total, err := h.service.ListUsers(page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = toResponse(&u)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": responses,
		"pagination": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.service.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created",
		"user":    toResponse(user),
	})
}

func getUserID(c echo.Context) int64 {
	if id, ok := c.Get("user_id").(int64); ok {
		return id
	}
	if id, ok := c.Get("user_id").(float64); ok {
		return int64(id)
	}
	return 0
}

func toResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		RealName:  user.RealName,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
