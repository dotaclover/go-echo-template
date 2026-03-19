package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService JWT 令牌服务
type JWTService struct {
	secret     string
	expiration time.Duration
}

// JWTClaims 自定义 Claims
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, expiration time.Duration) *JWTService {
	return &JWTService{secret: secret, expiration: expiration}
}

// Generate 生成 Token
func (s *JWTService) Generate(userID int64, username, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// Parse 解析并验证 Token
func (s *JWTService) Parse(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// Refresh 刷新 Token（返回新 Token）
func (s *JWTService) Refresh(tokenString string) (string, error) {
	claims, err := s.Parse(tokenString)
	if err != nil {
		return "", err
	}
	return s.Generate(claims.UserID, claims.Username, claims.Role)
}
