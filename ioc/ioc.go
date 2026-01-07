package ioc

import (
	"fmt"
	"sync"
)

// Provider 提供依赖的工厂函数
type Provider func() (any, error)

// Container IOC 容器
type Container struct {
	mu        sync.RWMutex
	data      map[string]interface{} // 存储实例
	providers map[string]Provider    // 存储 provider 函数
	singleton map[string]bool        // 标记是否为单例
	instances map[string]interface{} // 存储已创建的实例（用于单例）
	once      map[string]*sync.Once  // 用于单例的并发控制
}

var Default = New()

// New 创建新的容器实例
func New() *Container {
	return &Container{
		data:      make(map[string]interface{}),
		providers: make(map[string]Provider),
		singleton: make(map[string]bool),
		instances: make(map[string]interface{}),
		once:      make(map[string]*sync.Once),
	}
}

// Provide 注册 provider 函数
// singleton 为 true 时，provider 只会被调用一次，后续获取返回同一个实例
func (c *Container) Provide(key string, provider Provider, singleton bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providers[key] = provider
	c.singleton[key] = singleton
	// 如果之前有直接存储的实例，清除它（provider 优先）
	delete(c.data, key)
	if singleton {
		c.once[key] = &sync.Once{}
	} else {
		delete(c.once, key)
		delete(c.instances, key)
	}
}

// Get 获取依赖，如果不存在返回 nil 和 false
// 如果注册了 provider，会自动调用 provider 创建实例
func (c *Container) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	// 先检查直接存储的实例
	value, ok := c.data[key]
	if ok {
		c.mu.RUnlock()
		return value, true
	}

	// 检查是否有 provider
	provider, hasProvider := c.providers[key]
	isSingleton := c.singleton[key]
	c.mu.RUnlock()

	if !hasProvider {
		return nil, false
	}

	// 如果是单例，使用 sync.Once 确保只创建一次
	if isSingleton {
		c.mu.Lock()
		once, exists := c.once[key]
		if !exists {
			once = &sync.Once{}
			c.once[key] = once
		}
		c.mu.Unlock()

		var err error
		once.Do(func() {
			instance, e := provider()
			if e != nil {
				err = e
				return
			}
			c.mu.Lock()
			c.instances[key] = instance
			c.mu.Unlock()
		})

		if err != nil {
			return nil, false
		}

		c.mu.RLock()
		value = c.instances[key]
		c.mu.RUnlock()
		return value, true
	}

	// 非单例，每次调用 provider
	instance, err := provider()
	if err != nil {
		return nil, false
	}
	return instance, true
}

// MustGet 必须获取依赖，如果不存在会 panic
// 如果注册了 provider，会自动调用 provider 创建实例
func (c *Container) MustGet(key string) interface{} {
	value, ok := c.Get(key)
	if !ok {
		panic(fmt.Sprintf("ioc: key '%s' not found in container", key))
	}
	return value
}

// Delete 删除依赖（包括实例、provider 和单例标记）
func (c *Container) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	delete(c.providers, key)
	delete(c.singleton, key)
	delete(c.instances, key)
	delete(c.once, key)
}

// Clear 清空所有依赖（包括实例、provider 等）
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]interface{})
	c.providers = make(map[string]Provider)
	c.singleton = make(map[string]bool)
	c.instances = make(map[string]interface{})
	c.once = make(map[string]*sync.Once)
}

// Keys 获取所有键（包括实例和 provider）
func (c *Container) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keyMap := make(map[string]bool)
	for k := range c.data {
		keyMap[k] = true
	}
	for k := range c.providers {
		keyMap[k] = true
	}
	keys := make([]string, 0, len(keyMap))
	for k := range keyMap {
		keys = append(keys, k)
	}
	return keys
}

// Has 检查键是否存在（包括实例和 provider）
func (c *Container) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, hasData := c.data[key]
	_, hasProvider := c.providers[key]
	return hasData || hasProvider
}

// HasProvider 检查是否有注册的 provider
func (c *Container) HasProvider(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.providers[key]
	return ok
}

// HasInstance 检查是否有已创建的实例（包括直接存储的和通过 provider 创建的）
func (c *Container) HasInstance(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, hasData := c.data[key]
	_, hasInstance := c.instances[key]
	return hasData || hasInstance
}

// 全局函数，操作默认容器

// Provide 在默认容器中注册 provider
func Provide(key string, provider Provider, singleton bool) {
	Default.Provide(key, provider, singleton)
}

// Get 从默认容器获取依赖
func Get(key string) (interface{}, bool) {
	return Default.Get(key)
}

// MustGet 从默认容器必须获取依赖 如果类型不匹配会panic
func MustGet[T any](key string) T {
	v := Default.MustGet(key)
	var anyV any = v
	typedValue, ok := anyV.(T)
	if !ok {
		panic(fmt.Sprintf("ioc: key '%s' not found in container", key))
	}
	return typedValue
}
