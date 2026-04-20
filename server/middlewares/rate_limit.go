package middlewares

import (
	"myapp/common"
	"myapp/services"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimitConfig 限流中间件配置
type RateLimitConfig struct {
	Limiter services.RateLimiterInterface
	Limit   int                         // 窗口内最大请求数
	Window  time.Duration               // 窗口时长
	KeyFunc func(c echo.Context) string // 自定义 key 生成（默认用 IP）
	Message string                      // 超限时的错误消息
}

// RateLimit 限流中间件
//
// 用法：
//
//	limiter := services.NewMemoryRateLimiter()
//	e.Use(middlewares.RateLimit(middlewares.RateLimitConfig{
//	    Limiter: limiter,
//	    Limit:   100,
//	    Window:  time.Minute,
//	}))
func RateLimit(cfg RateLimitConfig) echo.MiddlewareFunc {
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = func(c echo.Context) string {
			return c.RealIP()
		}
	}
	if cfg.Message == "" {
		cfg.Message = "Too many requests"
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := cfg.KeyFunc(c)
			allowed, err := cfg.Limiter.Allow(key, cfg.Limit, cfg.Window)
			if err != nil {
				// 限流器故障时放行，不影响业务
				return next(c)
			}
			if !allowed {
				return common.Error(c, common.TooManyRequestsError(cfg.Message))
			}
			return next(c)
		}
	}
}
