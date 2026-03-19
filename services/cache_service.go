package services

import (
	"encoding/json"
	"sync"
	"time"
)

// ============================================================================
// 内存缓存实现
// ============================================================================

type memoryCacheEntry struct {
	value     string
	expiresAt time.Time
}

// MemoryCacheService 内存缓存（单机场景）
type MemoryCacheService struct {
	data map[string]*memoryCacheEntry
	mu   sync.RWMutex
}

func NewMemoryCacheService() *MemoryCacheService {
	s := &MemoryCacheService{data: make(map[string]*memoryCacheEntry)}
	go s.cleanup()
	return s
}

func (s *MemoryCacheService) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return "", nil
	}
	return entry.value, nil
}

func (s *MemoryCacheService) Set(key string, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = &memoryCacheEntry{value: value, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (s *MemoryCacheService) Del(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *MemoryCacheService) Exists(key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.data[key]
	return ok && time.Now().Before(entry.expiresAt), nil
}

func (s *MemoryCacheService) GetOrSet(key string, ttl time.Duration, fn func() (string, error)) (string, error) {
	val, _ := s.Get(key)
	if val != "" {
		return val, nil
	}
	result, err := fn()
	if err != nil {
		return "", err
	}
	_ = s.Set(key, result, ttl)
	return result, nil
}

func (s *MemoryCacheService) cleanup() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for k, v := range s.data {
			if now.After(v.expiresAt) {
				delete(s.data, k)
			}
		}
		s.mu.Unlock()
	}
}

// ============================================================================
// Redis 缓存实现
// ============================================================================

// RedisCacheService 基于 Redis 的缓存实现
type RedisCacheService struct {
	redis *RedisService
}

func NewRedisCacheService(redis *RedisService) *RedisCacheService {
	return &RedisCacheService{redis: redis}
}

func (s *RedisCacheService) Get(key string) (string, error) {
	return s.redis.Get(key)
}

func (s *RedisCacheService) Set(key string, value string, ttl time.Duration) error {
	return s.redis.Set(key, value, ttl)
}

func (s *RedisCacheService) Del(key string) error {
	return s.redis.Del(key)
}

func (s *RedisCacheService) Exists(key string) (bool, error) {
	return s.redis.Exists(key)
}

func (s *RedisCacheService) GetOrSet(key string, ttl time.Duration, fn func() (string, error)) (string, error) {
	val, err := s.redis.Get(key)
	if err != nil {
		return "", err
	}
	if val != "" {
		return val, nil
	}
	result, err := fn()
	if err != nil {
		return "", err
	}
	_ = s.redis.Set(key, result, ttl)
	return result, nil
}

// ============================================================================
// 缓存便捷方法（JSON 序列化）
// ============================================================================

// CacheGetJSON 从缓存获取并反序列化
func CacheGetJSON(cache CacheInterface, key string, dest interface{}) error {
	val, err := cache.Get(key)
	if err != nil || val == "" {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// CacheSetJSON 序列化后存入缓存
func CacheSetJSON(cache CacheInterface, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cache.Set(key, string(data), ttl)
}
