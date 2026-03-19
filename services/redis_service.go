package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService Redis 连接管理 + 常用操作
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// NewRedisService 创建 Redis 服务
func NewRedisService(cfg RedisConfig) (*RedisService, error) {
	addr := cfg.Host + ":" + cfg.Port
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 10
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connect failed: %w", err)
	}

	return &RedisService{client: client, ctx: ctx}, nil
}

// Client 获取原始 redis.Client（需要高级操作时使用）
func (r *RedisService) Client() *redis.Client {
	return r.client
}

// Close 关闭连接
func (r *RedisService) Close() error {
	return r.client.Close()
}

// ===== String 操作 =====

func (r *RedisService) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisService) Set(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisService) Del(keys ...string) error {
	return r.client.Del(r.ctx, keys...).Err()
}

func (r *RedisService) Exists(key string) (bool, error) {
	n, err := r.client.Exists(r.ctx, key).Result()
	return n > 0, err
}

func (r *RedisService) Expire(key string, ttl time.Duration) error {
	return r.client.Expire(r.ctx, key, ttl).Err()
}

func (r *RedisService) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

func (r *RedisService) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisService) IncrBy(key string, value int64) (int64, error) {
	return r.client.IncrBy(r.ctx, key, value).Result()
}

// ===== Hash 操作 =====

func (r *RedisService) HGet(key, field string) (string, error) {
	val, err := r.client.HGet(r.ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisService) HSet(key string, values ...interface{}) error {
	return r.client.HSet(r.ctx, key, values...).Err()
}

func (r *RedisService) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

func (r *RedisService) HDel(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// ===== List 操作 =====

func (r *RedisService) LPush(key string, values ...interface{}) error {
	return r.client.LPush(r.ctx, key, values...).Err()
}

func (r *RedisService) RPop(key string) (string, error) {
	val, err := r.client.RPop(r.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisService) LLen(key string) (int64, error) {
	return r.client.LLen(r.ctx, key).Result()
}

// ===== Set 操作 =====

func (r *RedisService) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

func (r *RedisService) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

func (r *RedisService) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// ===== JSON 便捷方法 =====

// SetJSON 序列化为 JSON 后存储
func (r *RedisService) SetJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(key, string(data), ttl)
}

// GetJSON 获取并反序列化 JSON
func (r *RedisService) GetJSON(key string, dest interface{}) error {
	val, err := r.Get(key)
	if err != nil {
		return err
	}
	if val == "" {
		return redis.Nil
	}
	return json.Unmarshal([]byte(val), dest)
}

// ===== Pub/Sub =====

func (r *RedisService) Publish(channel string, message interface{}) error {
	return r.client.Publish(r.ctx, channel, message).Err()
}

func (r *RedisService) Subscribe(channels ...string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channels...)
}
