package middlewares

import (
	"myapp/common"
	"myapp/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTAuth JWT 认证中间件
func JWTAuth(jwtService *services.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return common.Error(c, common.UnauthorizedError("Missing authorization header"))
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return common.Error(c, common.UnauthorizedError("Invalid authorization header format"))
			}

			claims, err := jwtService.Parse(parts[1])
			if err != nil {
				return common.Error(c, common.UnauthorizedError("Invalid token"))
			}

			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

// AdminOnly 管理员权限中间件（需先经过 JWTAuth）
func AdminOnly() echo.MiddlewareFunc {
	return RequireRoles("admin")
}

// RequireRoles 指定角色权限中间件（需先经过 JWTAuth）
func RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get("role").(string)
			for _, allowed := range roles {
				if role == allowed {
					return next(c)
				}
			}
			return common.Error(c, common.ForbiddenError(http.StatusText(http.StatusForbidden)))
		}
	}
}
