package cache

import (
	"sync"
	"time"

	kurohelperservice "kurohelperservice"
)

// Entry 單筆快取項目
type Entry[T any] struct {
	Value    T
	ExpireAt time.Time
}

// Store 泛型快取儲存
type Store[T any] struct {
	data       map[string]*Entry[T]
	expireTime time.Duration
	mu         sync.RWMutex
}

// NewStore 建立新的快取儲存
func NewStore[T any](expireTime time.Duration) *Store[T] {
	return &Store[T]{
		data:       make(map[string]*Entry[T]),
		expireTime: expireTime,
	}
}

// Set 設定快取
func (c *Store[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &Entry[T]{
		Value:    value,
		ExpireAt: time.Now().Add(c.expireTime),
	}
}

// Get 從快取中取得資料
func (c *Store[T]) Get(key string) (T, error) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()

	// 不存在或已過期
	if !ok || time.Now().After(item.ExpireAt) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		var zero T
		return zero, kurohelperservice.ErrCacheLost
	}

	return item.Value, nil
}

// Clean 清除過期快取
func (c *Store[T]) Clean() (deleteCount int, total int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	total = len(c.data)
	for k, d := range c.data {
		if time.Now().After(d.ExpireAt) {
			delete(c.data, k)
			deleteCount++
		}
	}

	return
}
