package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ============================================================================
// Redis 限流器（滑动窗口）
// ============================================================================

// RedisRateLimiter 基于 Redis 的滑动窗口限流器
type RedisRateLimiter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRateLimiter(redisSvc *RedisService) *RedisRateLimiter {
	return &RedisRateLimiter{client: redisSvc.Client(), ctx: context.Background()}
}

// Allow 判断是否允许请求
// key: 限流标识（如 "api:user:123"）
// limit: 窗口内允许的最大请求数
// window: 窗口时长
func (r *RedisRateLimiter) Allow(key string, limit int, window time.Duration) (bool, error) {
	rateKey := "rate:" + key
	now := time.Now().UnixMilli()
	windowStart := now - window.Milliseconds()

	pipe := r.client.Pipeline()
	// 移除窗口外的记录
	pipe.ZRemRangeByScore(r.ctx, rateKey, "0", fmt.Sprintf("%d", windowStart))
	// 统计窗口内的请求数
	countCmd := pipe.ZCard(r.ctx, rateKey)
	// 添加当前请求
	pipe.ZAdd(r.ctx, rateKey, redis.Z{Score: float64(now), Member: now})
	// 设置 key 过期（窗口 * 2，防止残留）
	pipe.Expire(r.ctx, rateKey, window*2)

	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return false, err
	}

	count := countCmd.Val()
	return count < int64(limit), nil
}

// ============================================================================
// 内存限流器（单机场景）
// ============================================================================

type memoryRateEntry struct {
	timestamps []int64
}

// MemoryRateLimiter 基于内存的滑动窗口限流器
type MemoryRateLimiter struct {
	data map[string]*memoryRateEntry
	mu   sync.Mutex
}

func NewMemoryRateLimiter() *MemoryRateLimiter {
	r := &MemoryRateLimiter{data: make(map[string]*memoryRateEntry)}
	go r.cleanup()
	return r
}

func (r *MemoryRateLimiter) Allow(key string, limit int, window time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UnixMilli()
	windowStart := now - window.Milliseconds()

	entry, exists := r.data[key]
	if !exists {
		entry = &memoryRateEntry{}
		r.data[key] = entry
	}

	// 移除窗口外的时间戳
	filtered := entry.timestamps[:0]
	for _, ts := range entry.timestamps {
		if ts > windowStart {
			filtered = append(filtered, ts)
		}
	}
	entry.timestamps = filtered

	if len(entry.timestamps) >= limit {
		return false, nil
	}

	entry.timestamps = append(entry.timestamps, now)
	return true, nil
}

func (r *MemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		r.mu.Lock()
		now := time.Now().UnixMilli()
		for k, v := range r.data {
			// 清理 10 分钟无活动的 key
			if len(v.timestamps) == 0 || v.timestamps[len(v.timestamps)-1] < now-600000 {
				delete(r.data, k)
			}
		}
		r.mu.Unlock()
	}
}
