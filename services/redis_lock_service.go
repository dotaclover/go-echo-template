package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisLockService Redis 分布式锁
type RedisLockService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisLockService(redisSvc *RedisService) *RedisLockService {
	return &RedisLockService{client: redisSvc.Client(), ctx: context.Background()}
}

func (s *RedisLockService) Obtain(key string, ttl time.Duration, retries int) (Lock, error) {
	lockKey := "lock:" + key
	value := uuid.New().String()

	for i := 0; i <= retries; i++ {
		ok, err := s.client.SetNX(s.ctx, lockKey, value, ttl).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			return &redisLock{client: s.client, ctx: s.ctx, key: lockKey, value: value}, nil
		}
		if i < retries {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil, errors.New("failed to obtain redis lock")
}

func (s *RedisLockService) TryObtain(key string, ttl time.Duration) (Lock, error) {
	return s.Obtain(key, ttl, 0)
}

func (s *RedisLockService) WithLock(key string, ttl time.Duration, fn func() error) error {
	lock, err := s.Obtain(key, ttl, 3)
	if err != nil {
		return err
	}
	defer lock.Release()
	return fn()
}

// redisLock Redis 锁实例
type redisLock struct {
	client *redis.Client
	ctx    context.Context
	key    string
	value  string
}

// Release 释放锁（Lua 脚本保证原子性）
func (l *redisLock) Release() error {
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		end
		return 0
	`
	result, err := l.client.Eval(l.ctx, script, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return errors.New("lock not owned")
	}
	return nil
}

// Refresh 续期锁
func (l *redisLock) Refresh(ttl time.Duration) error {
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		end
		return 0
	`
	result, err := l.client.Eval(l.ctx, script, []string{l.key}, l.value, ttl.Milliseconds()).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return errors.New("lock not owned")
	}
	return nil
}
