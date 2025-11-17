// Package helper -----------------------------
// @file      : cache.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/1/16
// -------------------------------------------
package helper

import (
	"sync"
	"time"
)

// Cache 通用的带过期时间的缓存结构
type Cache struct {
	data      map[string]interface{} // 缓存数据
	expiresAt map[string]time.Time   // 过期时间
	mu        sync.RWMutex           // 读写锁
	ttl       time.Duration          // 默认缓存过期时间
}

// NewCache 创建一个新的缓存实例
// ttl: 缓存过期时间，如果为0则使用默认值
func NewCache(defaultTTL time.Duration) *Cache {
	if defaultTTL <= 0 {
		defaultTTL = 30 * time.Second // 默认30秒
	}
	return &Cache{
		data:      make(map[string]interface{}),
		expiresAt: make(map[string]time.Time),
		ttl:       defaultTTL,
	}
}

// Get 从缓存获取数据，如果过期或不存在则返回 false
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	value, exists := c.data[key]
	expireTime, hasExpire := c.expiresAt[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if !hasExpire || time.Now().After(expireTime) {
		// 缓存已过期，删除
		c.mu.Lock()
		// 双重检查，防止并发删除
		if et, ok := c.expiresAt[key]; ok && (time.Now().After(et) || !hasExpire) {
			delete(c.data, key)
			delete(c.expiresAt, key)
		}
		c.mu.Unlock()
		return nil, false
	}

	return value, true
}

// Set 设置缓存数据，使用默认 TTL
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.ttl)
}

// SetWithTTL 设置缓存数据，使用自定义 TTL
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	if ttl <= 0 {
		ttl = c.ttl
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiresAt[key] = time.Now().Add(ttl)
}

// Delete 删除缓存项
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	delete(c.expiresAt, key)
}

// Clear 清空所有缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]interface{})
	c.expiresAt = make(map[string]time.Time)
}

// CleanExpired 清理过期缓存
func (c *Cache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, expireTime := range c.expiresAt {
		if now.After(expireTime) {
			delete(c.data, key)
			delete(c.expiresAt, key)
		}
	}
}

// Size 返回当前缓存项数量
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// GetStringSlice 获取字符串切片类型的缓存值
func (c *Cache) GetStringSlice(key string) ([]string, bool) {
	value, ok := c.Get(key)
	if !ok {
		return nil, false
	}

	if slice, ok := value.([]string); ok {
		return slice, true
	}
	return nil, false
}

// SetStringSlice 设置字符串切片类型的缓存值
func (c *Cache) SetStringSlice(key string, value []string) {
	c.Set(key, value)
}

// SetStringSliceWithTTL 设置字符串切片类型的缓存值，使用自定义 TTL
func (c *Cache) SetStringSliceWithTTL(key string, value []string, ttl time.Duration) {
	c.SetWithTTL(key, value, ttl)
}
