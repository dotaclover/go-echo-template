package user

import (
	"myapp/common"
	"myapp/models"
	"myapp/services"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Service 用户服务接口
type Service interface {
	Register(req *RegisterRequest) (*models.User, error)
	CreateUser(req *CreateUserRequest) (*models.User, error)
	Login(req *LoginRequest) (*services.TokenPair, *models.User, error)
	RefreshToken(refreshToken string) (*services.TokenPair, error)
	GetProfile(userID int64) (*models.User, error)
	UpdateProfile(userID int64, req *UpdateProfileRequest) (*models.User, error)
	UpdatePassword(userID int64, req *UpdatePasswordRequest) error
	GetUserByID(id int64) (*models.User, error)
	UpdateUser(id int64, req *UpdateUserRequest) (*models.User, error)
	UpdateUserStatus(id int64, status string, operatorID int64) (*models.User, error)
	DeleteUser(id int64, operatorID int64) error
	ListUsers(page, pageSize int) ([]models.User, int64, error)
}

type service struct {
	repo       Repository
	jwtService *services.JWTService
}

func NewService(repo Repository, jwtService *services.JWTService) Service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) Register(req *RegisterRequest) (*models.User, error) {
	return s.createUser(&CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		RealName: req.RealName,
		Role:     models.RoleUser,
		Status:   models.StatusActive,
	})
}

func (s *service) CreateUser(req *CreateUserRequest) (*models.User, error) {
	return s.createUser(req)
}

func (s *service) Login(req *LoginRequest) (*services.TokenPair, *models.User, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, common.UnauthorizedError("invalid username or password")
	}
	if user.Status != models.StatusActive {
		return nil, nil, common.ForbiddenError("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, common.UnauthorizedError("invalid username or password")
	}

	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, nil, common.InternalError("failed to generate token")
	}
	return tokens, user, nil
}

func (s *service) RefreshToken(refreshToken string) (*services.TokenPair, error) {
	tokens, err := s.jwtService.Refresh(refreshToken)
	if err != nil {
		return nil, common.UnauthorizedError("invalid refresh token")
	}
	return tokens, nil
}

func (s *service) GetProfile(userID int64) (*models.User, error) {
	return s.repo.FindByID(userID)
}

func (s *service) UpdateProfile(userID int64, req *UpdateProfileRequest) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.RealName != nil {
		user.RealName = strings.TrimSpace(*req.RealName)
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email != "" {
			existing, err := s.repo.FindByEmail(email)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != user.ID {
				return nil, common.ConflictError("email already exists")
			}
		}
		user.Email = email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) UpdatePassword(userID int64, req *UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return common.BadRequestError("old password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return common.InternalError("failed to hash password")
	}

	user.PasswordHash = string(hash)
	return s.repo.Update(user)
}

func (s *service) GetUserByID(id int64) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) UpdateUser(id int64, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email != "" {
			existing, err := s.repo.FindByEmail(email)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != user.ID {
				return nil, common.ConflictError("email already exists")
			}
		}
		user.Email = email
	}
	if req.RealName != nil {
		user.RealName = strings.TrimSpace(*req.RealName)
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateUserStatus(id int64, status string, operatorID int64) (*models.User, error) {
	if id == operatorID {
		return nil, common.BadRequestError("cannot change your own status")
	}
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	user.Status = status
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) DeleteUser(id int64, operatorID int64) error {
	if id == operatorID {
		return common.BadRequestError("cannot delete your own account")
	}
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *service) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(page, pageSize)
}

func (s *service) createUser(req *CreateUserRequest) (*models.User, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, common.ValidationError("validation failed", map[string]interface{}{"username": "username is required"})
	}
	existing, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, common.ConflictError("username already exists")
	}
	email := strings.TrimSpace(req.Email)
	if email != "" {
		existingByEmail, err := s.repo.FindByEmail(email)
		if err != nil {
			return nil, err
		}
		if existingByEmail != nil {
			return nil, common.ConflictError("email already exists")
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.InternalError("failed to hash password")
	}
	status := req.Status
	if status == "" {
		status = models.StatusActive
	}
	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
		Email:        email,
		RealName:     strings.TrimSpace(req.RealName),
		Role:         req.Role,
		Status:       status,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
