package user

import (
	"myapp/common"
	"myapp/models"
	"myapp/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Handler 用户处理器
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(g *echo.Group) {
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.RefreshToken)
}

func (h *Handler) RegisterAuthRoutes(g *echo.Group) {
	g.GET("/profile", h.GetProfile)
	g.PUT("/profile", h.UpdateProfile)
	g.POST("/password", h.UpdatePassword)
}

func (h *Handler) RegisterAdminRoutes(g *echo.Group) {
	g.GET("", h.ListUsers)
	g.GET("/:id", h.GetUser)
	g.POST("", h.CreateUser)
	g.PUT("/:id", h.UpdateUser)
	g.PATCH("/:id/status", h.UpdateUserStatus)
	g.DELETE("/:id", h.DeleteUser)
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}
	user, err := h.service.Register(&req)
	if err != nil {
		return common.Error(c, err)
	}
	return common.Created(c, "User registered", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}

	tokens, user, err := h.service.Login(&req)
	if err != nil {
		return common.Error(c, err)
	}

	return common.Success(c, "Login successful", AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
		User:         toResponse(user),
	})
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}
	tokens, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return common.Error(c, err)
	}
	return common.Success(c, "Token refreshed", map[string]interface{}{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"token_type":    tokens.TokenType,
		"expires_in":    tokens.ExpiresIn,
	})
}

func (h *Handler) GetProfile(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return common.Error(c, common.UnauthorizedError("Unauthorized"))
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		return common.Error(c, err)
	}

	return common.Success(c, "Profile fetched", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return common.Error(c, common.UnauthorizedError("Unauthorized"))
	}

	var req UpdateProfileRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}

	user, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		return common.Error(c, err)
	}

	return common.Success(c, "Profile updated", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) UpdatePassword(c echo.Context) error {
	userID := getUserID(c)
	if userID == 0 {
		return common.Error(c, common.UnauthorizedError("Unauthorized"))
	}

	var req UpdatePasswordRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}

	if err := h.service.UpdatePassword(userID, &req); err != nil {
		return common.Error(c, err)
	}

	return common.Success(c, "Password updated", nil)
}

func (h *Handler) ListUsers(c echo.Context) error {
	page, pageSize, err := utils.Validator.ValidatePagination(c.QueryParam("page"), c.QueryParam("page_size"))
	if err != nil {
		return common.Error(c, common.BadRequestError(err.Error()))
	}

	users, total, err := h.service.ListUsers(page, pageSize)
	if err != nil {
		return common.Error(c, err)
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = toResponse(&u)
	}

	totalPages := int64(0)
	if pageSize > 0 {
		totalPages = (total + int64(pageSize) - 1) / int64(pageSize)
	}
	return common.Paginated(c, "Users fetched", map[string]interface{}{"users": responses}, &common.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *Handler) GetUser(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return common.Error(c, common.BadRequestError("invalid user id"))
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return common.Error(c, err)
	}
	return common.Success(c, "User fetched", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		return common.Error(c, err)
	}

	return common.Created(c, "User created", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return common.Error(c, common.BadRequestError("invalid user id"))
	}
	var req UpdateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}
	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		return common.Error(c, err)
	}
	return common.Success(c, "User updated", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) UpdateUserStatus(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return common.Error(c, common.BadRequestError("invalid user id"))
	}
	var req UpdateUserStatusRequest
	if err := bindAndValidate(c, &req); err != nil {
		return common.Error(c, err)
	}
	user, err := h.service.UpdateUserStatus(id, req.Status, getUserID(c))
	if err != nil {
		return common.Error(c, err)
	}
	return common.Success(c, "User status updated", map[string]interface{}{"user": toResponse(user)})
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return common.Error(c, common.BadRequestError("invalid user id"))
	}
	if err := h.service.DeleteUser(id, getUserID(c)); err != nil {
		return common.Error(c, err)
	}
	return common.Success(c, "User deleted", nil)
}

func bindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return common.BadRequestError("invalid request body")
	}
	if err := c.Validate(req); err != nil {
		return common.ValidationError("validation failed", map[string]interface{}{"fields": utils.FormatValidationErrors(err)})
	}
	return nil
}

func parseID(id string) (int64, error) {
	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, strconv.ErrSyntax
	}
	return parsed, nil
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
