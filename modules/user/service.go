package user

import (
	"errors"
	"myapp/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service 用户服务接口
type Service interface {
	Register(req *RegisterRequest) (*models.User, error)
	Login(req *LoginRequest) (string, *models.User, error)
	GetProfile(userID int64) (*models.User, error)
	UpdateProfile(userID int64, req *UpdateProfileRequest) (*models.User, error)
	UpdatePassword(userID int64, req *UpdatePasswordRequest) error
	ListUsers(page, pageSize int) ([]models.User, int64, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) Register(req *RegisterRequest) (*models.User, error) {
	existing, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		Email:        req.Email,
		RealName:     req.RealName,
		Role:         models.RoleUser,
		Status:       models.StatusActive,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) Login(req *LoginRequest) (string, *models.User, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid username or password")
	}
	if user.Status != models.StatusActive {
		return "", nil, errors.New("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
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
		user.RealName = *req.RealName
	}
	if req.Email != nil {
		user.Email = *req.Email
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
		return errors.New("old password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.repo.Update(user)
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

func (s *service) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
