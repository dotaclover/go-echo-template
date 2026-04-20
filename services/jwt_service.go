package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService JWT 令牌服务
type JWTService struct {
	secret            string
	expiration        time.Duration
	refreshExpiration time.Duration
}

// JWTClaims 自定义 Claims
type JWTClaims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewJWTService(secret string, expiration, refreshExpiration time.Duration) *JWTService {
	return &JWTService{secret: secret, expiration: expiration, refreshExpiration: refreshExpiration}
}

// Generate 生成 access token
func (s *JWTService) Generate(userID int64, username, role string) (string, error) {
	return s.generate(userID, username, role, "access", s.expiration)
}

// GenerateRefresh 生成 refresh token
func (s *JWTService) GenerateRefresh(userID int64, username, role string) (string, error) {
	return s.generate(userID, username, role, "refresh", s.refreshExpiration)
}

// GenerateTokenPair 生成 access + refresh token
func (s *JWTService) GenerateTokenPair(userID int64, username, role string) (*TokenPair, error) {
	accessToken, err := s.Generate(userID, username, role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.GenerateRefresh(userID, username, role)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.expiration.Seconds()),
	}, nil
}

func (s *JWTService) generate(userID int64, username, role, tokenType string, expiration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
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
func (s *JWTService) Refresh(tokenString string) (*TokenPair, error) {
	claims, err := s.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}
	return s.GenerateTokenPair(claims.UserID, claims.Username, claims.Role)
}
