package user

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	RealName string `json:"real_name" validate:"omitempty,max=50"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	RealName *string `json:"real_name" validate:"omitempty,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

// CreateUserRequest 管理员创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	RealName string `json:"real_name" validate:"omitempty,max=50"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	Status   string `json:"status" validate:"omitempty,oneof=active disabled"`
}

// UpdateUserRequest 管理员更新用户请求
type UpdateUserRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	RealName *string `json:"real_name" validate:"omitempty,max=50"`
	Role     *string `json:"role" validate:"omitempty,oneof=admin user"`
}

// UpdateUserStatusRequest 管理员更新用户状态请求
type UpdateUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active disabled"`
}

// AuthResponse 登录/刷新响应
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RealName  string    `json:"real_name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
