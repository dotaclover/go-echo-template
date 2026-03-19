package router

import (
	"myapp/config"
	"myapp/modules/user"
	"myapp/server/database"
	"myapp/server/middlewares"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	// 设置 JWT 密钥
	middlewares.SetJWTSecret(cfg.JWT.Secret)

	// ===== 健康检查（公开）=====
	e.GET("/health", func(c echo.Context) error {
		if err := database.HealthCheck(db); err != nil {
			return c.JSON(503, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	// ===== 用户模块（自带路由注册）=====
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg.JWT.Secret)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(e)

	// ===== API 路由组 =====
	apiGroup := e.Group("/api")

	// 需要认证的路由
	authGroup := apiGroup.Group("")
	authGroup.Use(middlewares.JWTAuth())

	// 管理员路由
	adminGroup := apiGroup.Group("")
	adminGroup.Use(middlewares.JWTAuth(), middlewares.AdminOnly())

	// ===== 在此注册业务模块路由 =====
	// example:
	// product.RegisterRoutes(authGroup, db)
	// order.RegisterAdminRoutes(adminGroup, db)

	_ = authGroup
	_ = adminGroup
}
