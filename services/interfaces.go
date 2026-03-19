package services

import "time"

// ============================================================================
// 缓存接口
// ============================================================================

// CacheInterface 缓存服务接口
type CacheInterface interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Del(key string) error
	Exists(key string) (bool, error)
	GetOrSet(key string, ttl time.Duration, fn func() (string, error)) (string, error)
}

// ============================================================================
// 锁接口
// ============================================================================

// LockInterface 锁服务接口
type LockInterface interface {
	Obtain(key string, ttl time.Duration, retries int) (Lock, error)
	TryObtain(key string, ttl time.Duration) (Lock, error)
	WithLock(key string, ttl time.Duration, fn func() error) error
}

// Lock 锁实例接口
type Lock interface {
	Release() error
	Refresh(ttl time.Duration) error
}

// ============================================================================
// 限流接口
// ============================================================================

// RateLimiterInterface 限流器接口
type RateLimiterInterface interface {
	Allow(key string, limit int, window time.Duration) (bool, error)
}

// ============================================================================
// 通知接口
// ============================================================================

// Notifier 通知渠道接口
type Notifier interface {
	Name() string
	Send(title, content string) error
}

// ============================================================================
// 短信接口
// ============================================================================

// SMSSender 短信发送接口
type SMSSender interface {
	Send(phone, templateID string, params map[string]string) error
}
