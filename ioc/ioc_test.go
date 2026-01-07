package ioc

import (
	"errors"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() returned nil")
	}
	if c.data == nil {
		t.Fatal("data map is nil")
	}
	if c.providers == nil {
		t.Fatal("providers map is nil")
	}
}

func TestContainer_Provide(t *testing.T) {
	c := New()

	// 测试注册单例 provider
	c.Provide("singleton", func() (any, error) {
		return "singleton-value", nil
	}, true)

	if !c.HasProvider("singleton") {
		t.Fatal("provider not registered")
	}
	if !c.singleton["singleton"] {
		t.Fatal("singleton flag not set")
	}

	// 测试注册非单例 provider
	c.Provide("non-singleton", func() (any, error) {
		return "non-singleton-value", nil
	}, false)

	if c.singleton["non-singleton"] {
		t.Fatal("singleton flag should be false")
	}
}

func TestContainer_Get(t *testing.T) {
	c := New()

	// 测试获取不存在的 key
	_, ok := c.Get("not-exists")
	if ok {
		t.Fatal("should return false for non-existent key")
	}

	// 测试通过 provider 获取（单例）
	c.Provide("test", func() (any, error) {
		return "test-value", nil
	}, true)

	value, ok := c.Get("test")
	if !ok {
		t.Fatal("should return true for existing provider")
	}
	if value != "test-value" {
		t.Fatalf("expected 'test-value', got %v", value)
	}

	// 测试单例：多次获取应该是同一个实例
	value2, _ := c.Get("test")
	if value != value2 {
		t.Fatal("singleton should return same instance")
	}

	// 测试非单例：每次获取应该是新实例
	c.Provide("non-singleton", func() (any, error) {
		return &struct{ ID int }{ID: 1}, nil
	}, false)

	v1, _ := c.Get("non-singleton")
	v2, _ := c.Get("non-singleton")
	if v1 == v2 {
		t.Fatal("non-singleton should return different instances")
	}
}

func TestContainer_Get_ProviderError(t *testing.T) {
	c := New()

	c.Provide("error", func() (any, error) {
		return nil, errors.New("provider error")
	}, false)

	_, ok := c.Get("error")
	if ok {
		t.Fatal("should return false when provider returns error")
	}
}

func TestContainer_MustGet(t *testing.T) {
	c := New()

	// 测试获取存在的 key
	c.Provide("exists", func() (any, error) {
		return "exists-value", nil
	}, true)

	value := c.MustGet("exists")
	if value != "exists-value" {
		t.Fatalf("expected 'exists-value', got %v", value)
	}

	// 测试获取不存在的 key（应该 panic）
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustGet should panic for non-existent key")
		}
	}()
	c.MustGet("not-exists")
}

func TestContainer_Delete(t *testing.T) {
	c := New()

	c.Provide("to-delete", func() (any, error) {
		return "value", nil
	}, true)

	// 先获取一次，确保实例被创建
	c.Get("to-delete")

	// 删除
	c.Delete("to-delete")

	// 验证已删除
	if c.Has("to-delete") {
		t.Fatal("key should be deleted")
	}
	if c.HasProvider("to-delete") {
		t.Fatal("provider should be deleted")
	}
	if c.HasInstance("to-delete") {
		t.Fatal("instance should be deleted")
	}
}

func TestContainer_Clear(t *testing.T) {
	c := New()

	c.Provide("key1", func() (any, error) { return "v1", nil }, true)
	c.Provide("key2", func() (any, error) { return "v2", nil }, false)

	// 获取一次，创建实例
	c.Get("key1")
	c.Get("key2")

	c.Clear()

	if len(c.Keys()) != 0 {
		t.Fatal("container should be empty after Clear()")
	}
}

func TestContainer_Has(t *testing.T) {
	c := New()

	if c.Has("not-exists") {
		t.Fatal("should return false for non-existent key")
	}

	c.Provide("exists", func() (any, error) {
		return "value", nil
	}, true)

	if !c.Has("exists") {
		t.Fatal("should return true for existing provider")
	}
}

func TestContainer_HasProvider(t *testing.T) {
	c := New()

	if c.HasProvider("not-exists") {
		t.Fatal("should return false for non-existent provider")
	}

	c.Provide("exists", func() (any, error) {
		return "value", nil
	}, true)

	if !c.HasProvider("exists") {
		t.Fatal("should return true for existing provider")
	}
}

func TestContainer_HasInstance(t *testing.T) {
	c := New()

	if c.HasInstance("not-exists") {
		t.Fatal("should return false for non-existent instance")
	}

	c.Provide("exists", func() (any, error) {
		return "value", nil
	}, true)

	// 获取前应该没有实例
	if c.HasInstance("exists") {
		t.Fatal("should return false before Get()")
	}

	// 获取后应该有实例
	c.Get("exists")
	if !c.HasInstance("exists") {
		t.Fatal("should return true after Get()")
	}
}

func TestContainer_Keys(t *testing.T) {
	c := New()

	keys := c.Keys()
	if len(keys) != 0 {
		t.Fatal("empty container should return empty keys")
	}

	c.Provide("key1", func() (any, error) { return "v1", nil }, true)
	c.Provide("key2", func() (any, error) { return "v2", nil }, false)

	keys = c.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestContainer_ConcurrentAccess(t *testing.T) {
	c := New()
	c.Provide("test", func() (any, error) {
		return "value", nil
	}, true)

	var wg sync.WaitGroup
	goroutines := 100

	// 并发获取
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = c.Get("test")
		}()
	}

	// 并发注册
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := "key" + string(rune(id))
			c.Provide(key, func() (any, error) {
				return id, nil
			}, false)
		}(i)
	}

	wg.Wait()

	// 验证单例仍然正确
	value1, _ := c.Get("test")
	value2, _ := c.Get("test")
	if value1 != value2 {
		t.Fatal("singleton should be thread-safe")
	}
}

func TestDefault(t *testing.T) {
	if Default == nil {
		t.Fatal("Default container should not be nil")
	}

	// 测试默认容器
	Provide("default-test", func() (any, error) {
		return "default-value", nil
	}, true)

	value, ok := Get("default-test")
	if !ok {
		t.Fatal("should get value from default container")
	}
	if value != "default-value" {
		t.Fatalf("expected 'default-value', got %v", value)
	}

	// 清理
	Default.Delete("default-test")
}

func TestMustGet_Generic(t *testing.T) {
	Provide("string-value", func() (any, error) {
		return "test-string", nil
	}, true)

	value := MustGet[string]("string-value")
	if value != "test-string" {
		t.Fatalf("expected 'test-string', got %v", value)
	}

	// 清理
	Default.Delete("string-value")
}
