package user

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	RealName string `json:"real_name" validate:"omitempty,max=50"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	RealName *string `json:"real_name" validate:"omitempty,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
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
