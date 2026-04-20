package commands

import (
	"context"
	"fmt"
	"myapp/common"
	"myapp/config"
	"myapp/server/database"
	"myapp/server/middlewares"
	"myapp/server/router"
	"myapp/services"
	"myapp/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ServeCommand 启动服务命令
type ServeCommand struct{}

func (c *ServeCommand) Name() string        { return "serve" }
func (c *ServeCommand) Description() string { return "Start HTTP server" }

// CustomValidator Echo 验证器
type CustomValidator struct {
	validator *utils.ValidationService
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.ValidateStruct(i)
}

func (c *ServeCommand) Execute(args []string) error {
	flags, _ := ParseFlags(args)

	_ = godotenv.Load()
	utils.InitLogger()

	cfg := config.Load()
	if err := cfg.IsValid(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	if cfg.JWT.Secret == config.DefaultJWTSecret {
		utils.Logger.Warn("[SECURITY] JWT secret is still default, please set JWT_SECRET")
	}

	// 从命令行参数或配置获取主机和端口
	host := flags["host"]
	if host == "" {
		host = cfg.App.Host
	}
	port := flags["port"]
	if port == "" {
		port = cfg.App.Port
	}

	// 确保数据目录存在
	_ = os.MkdirAll("./data", 0755)

	// 初始化数据库
	db := database.InitDB()
	if err := database.HealthCheck(db); err != nil {
		database.CloseDB(db)
		return fmt.Errorf("database health check failed: %v", err)
	}

	// 创建 Echo 实例
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &CustomValidator{validator: utils.Validator}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		_ = common.Error(c, err)
	}

	e.Use(middlewares.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	if cfg.App.RateLimitEnabled {
		limiter := services.NewMemoryRateLimiter()
		e.Use(middlewares.RateLimit(middlewares.RateLimitConfig{
			Limiter: limiter,
			Limit:   cfg.App.RateLimitLimit,
			Window:  cfg.App.RateLimitWindow,
		}))
	}
	if cfg.App.Debug {
		e.Use(middleware.Logger())
	}

	// 注册路由
	router.RegisterRoutes(e, db, cfg)

	addr := host + ":" + port
	if host == "" {
		addr = ":" + port
	}
	utils.Logger.Infof("Server starting on %s", addr)

	serverErrCh := make(chan error, 1)

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		utils.Logger.Info("Shutting down...")
		database.CloseDB(db)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			utils.Logger.Errorf("Server shutdown error: %v", err)
		}
		utils.Logger.Info("Server stopped")
		os.Exit(0)
	}()

	go func() {
		server := &http.Server{
			Addr:         addr,
			ReadTimeout:  cfg.App.ReadTimeout,
			WriteTimeout: cfg.App.WriteTimeout,
			IdleTimeout:  cfg.App.IdleTimeout,
		}
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	if err := <-serverErrCh; err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	return nil
}
