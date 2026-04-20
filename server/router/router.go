package router

import (
	"myapp/common"
	"myapp/config"
	"myapp/modules/user"
	"myapp/server/database"
	"myapp/server/middlewares"
	"myapp/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	jwtService := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration, cfg.JWT.RefreshExpiration)

	// ===== 健康检查（公开）=====
	e.GET("/health/live", func(c echo.Context) error {
		return common.Success(c, "Service is live", map[string]string{"status": "live"})
	})
	e.GET("/health/ready", func(c echo.Context) error {
		if err := database.HealthCheck(db); err != nil {
			return common.Error(c, common.InternalError(err.Error()))
		}
		return common.Success(c, "Service is ready", map[string]string{"status": "ready"})
	})

	// ===== 用户模块 =====
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, jwtService)
	userHandler := user.NewHandler(userService)

	// ===== API 路由组 =====
	apiGroup := e.Group("/api/v1")
	publicGroup := apiGroup.Group("")

	// 需要认证的路由
	authGroup := apiGroup.Group("")
	authGroup.Use(middlewares.JWTAuth(jwtService))

	// 管理员路由
	adminGroup := apiGroup.Group("")
	adminGroup.Use(middlewares.JWTAuth(jwtService), middlewares.AdminOnly())

	authPublicGroup := publicGroup.Group("/auth")
	authPrivateGroup := authGroup.Group("/auth")
	adminUsersGroup := adminGroup.Group("/admin/users")

	userHandler.RegisterPublicRoutes(authPublicGroup)
	userHandler.RegisterAuthRoutes(authPrivateGroup)
	userHandler.RegisterAdminRoutes(adminUsersGroup)
}
