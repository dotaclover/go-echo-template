package services

import (
	"errors"
	"sync"
	"time"
)

// MemoryLockService 内存锁服务（单机场景）
type MemoryLockService struct {
	locks map[string]*lockEntry
	mu    sync.Mutex
}

type lockEntry struct {
	value     string
	expiresAt time.Time
}

func NewMemoryLockService() *MemoryLockService {
	s := &MemoryLockService{locks: make(map[string]*lockEntry)}
	go s.cleanup()
	return s
}

func (s *MemoryLockService) Obtain(key string, ttl time.Duration, retries int) (Lock, error) {
	for i := 0; i <= retries; i++ {
		s.mu.Lock()
		entry, exists := s.locks[key]
		if !exists || time.Now().After(entry.expiresAt) {
			value := time.Now().Format(time.RFC3339Nano)
			s.locks[key] = &lockEntry{value: value, expiresAt: time.Now().Add(ttl)}
			s.mu.Unlock()
			return &memoryLock{key: key, value: value, manager: s}, nil
		}
		s.mu.Unlock()
		if i < retries {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil, errors.New("failed to obtain lock")
}

func (s *MemoryLockService) TryObtain(key string, ttl time.Duration) (Lock, error) {
	return s.Obtain(key, ttl, 0)
}

func (s *MemoryLockService) WithLock(key string, ttl time.Duration, fn func() error) error {
	lock, err := s.Obtain(key, ttl, 3)
	if err != nil {
		return err
	}
	defer lock.Release()
	return fn()
}

func (s *MemoryLockService) cleanup() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for k, v := range s.locks {
			if now.After(v.expiresAt) {
				delete(s.locks, k)
			}
		}
		s.mu.Unlock()
	}
}

// memoryLock 内存锁实例
type memoryLock struct {
	key     string
	value   string
	manager *MemoryLockService
}

func (l *memoryLock) Release() error {
	l.manager.mu.Lock()
	defer l.manager.mu.Unlock()
	entry, exists := l.manager.locks[l.key]
	if !exists || entry.value != l.value {
		return errors.New("lock not owned")
	}
	delete(l.manager.locks, l.key)
	return nil
}

func (l *memoryLock) Refresh(ttl time.Duration) error {
	l.manager.mu.Lock()
	defer l.manager.mu.Unlock()
	entry, exists := l.manager.locks[l.key]
	if !exists || entry.value != l.value {
		return errors.New("lock not owned")
	}
	entry.expiresAt = time.Now().Add(ttl)
	return nil
}
