package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret string

// SetJWTSecret 设置 JWT 密钥（启动时调用）
func SetJWTSecret(secret string) {
	jwtSecret = secret
}

// JWTAuth JWT 认证中间件
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Missing authorization header",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid authorization header format",
				})
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid token claims",
				})
			}

			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", int64(userID))
			}
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}

			return next(c)
		}
	}
}

// AdminOnly 管理员权限中间件（需先经过 JWTAuth）
func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("role") != "admin" {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Admin permission required",
				})
			}
			return next(c)
		}
	}
}
