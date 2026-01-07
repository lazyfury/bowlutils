package eventbus

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Fatal("New() returned nil")
	}
	if bus.subs == nil {
		t.Fatal("subs map is nil")
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	bus := New()

	id, ch := bus.Subscribe("test-topic", 10)
	if id == 0 {
		t.Fatal("Subscribe() should return non-zero ID")
	}
	if ch == nil {
		t.Fatal("Subscribe() should return channel")
	}

	// 测试多次订阅
	id2, ch2 := bus.Subscribe("test-topic", 10)
	if id == id2 {
		t.Fatal("Subscribe() should return different IDs")
	}
	if ch == ch2 {
		t.Fatal("Subscribe() should return different channels")
	}
}

func TestEventBus_Subscribe_BufferSize(t *testing.T) {
	bus := New()

	// 测试 buffer <= 0 时应该使用 DefaultBufferSize (10)
	_, ch := bus.Subscribe("test", 0)
	if cap(ch) != DefaultBufferSize {
		t.Errorf("expected buffer size %d, got %d", DefaultBufferSize, cap(ch))
	}

	_, ch2 := bus.Subscribe("test2", -1)
	if cap(ch2) != DefaultBufferSize {
		t.Errorf("expected buffer size %d, got %d", DefaultBufferSize, cap(ch2))
	}

	// 测试正常 buffer
	_, ch3 := bus.Subscribe("test3", 5)
	if cap(ch3) != 5 {
		t.Errorf("expected buffer size 5, got %d", cap(ch3))
	}
}

func TestEventBus_Publish(t *testing.T) {
	bus := New()

	_, ch := bus.Subscribe("test-topic", 10)

	// 发布消息
	bus.Publish("test-topic", "test-message")

	// 接收消息
	select {
	case msg := <-ch:
		if msg != "test-message" {
			t.Errorf("expected 'test-message', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}
}

func TestEventBus_Publish_MultipleSubscribers(t *testing.T) {
	bus := New()

	_, ch1 := bus.Subscribe("test-topic", 10)
	_, ch2 := bus.Subscribe("test-topic", 10)

	bus.Publish("test-topic", "test-message")

	// 两个订阅者都应该收到消息
	select {
	case msg := <-ch1:
		if msg != "test-message" {
			t.Errorf("ch1: expected 'test-message', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("ch1: message not received")
	}

	select {
	case msg := <-ch2:
		if msg != "test-message" {
			t.Errorf("ch2: expected 'test-message', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("ch2: message not received")
	}
}

func TestEventBus_Publish_NoSubscribers(t *testing.T) {
	bus := New()

	// 发布到没有订阅者的主题（不应该 panic）
	bus.Publish("no-subscribers", "message")
}

func TestEventBus_Publish_FullChannel(t *testing.T) {
	bus := New()

	id, ch := bus.Subscribe("test-topic", 1) // 小 buffer

	// 填满 channel（通过发布消息）
	bus.Publish("test-topic", "blocking")

	// 发布消息（channel 已满，应该被丢弃）
	bus.Publish("test-topic", "should-be-dropped-1")

	// 再次发布（channel 已满，应该被丢弃）
	bus.Publish("test-topic", "should-be-dropped-2")

	// 验证丢弃计数（应该有 2 个被丢弃）
	dropped := bus.DroppedCount()
	if dropped != 2 {
		t.Errorf("expected 2 dropped messages, got %d", dropped)
	}

	// 接收第一个消息（channel 现在空了）
	select {
	case msg := <-ch:
		if msg != "blocking" {
			t.Errorf("expected 'blocking', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}

	// 再次发布（现在 channel 空了，应该能接收）
	bus.Publish("test-topic", "should-receive")

	select {
	case msg := <-ch:
		if msg != "should-receive" {
			t.Errorf("expected 'should-receive', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}

	// 清理
	bus.Unsubscribe("test-topic", id)
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := New()

	id, ch := bus.Subscribe("test-topic", 10)

	// 取消订阅
	bus.Unsubscribe("test-topic", id)

	// channel 应该被关闭
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("channel should be closed after Unsubscribe")
		}
	default:
		t.Fatal("channel should be closed")
	}

	// 发布消息，不应该发送到已取消的订阅者
	bus.Publish("test-topic", "message")
}

func TestEventBus_Unsubscribe_NonExistent(t *testing.T) {
	bus := New()

	// 取消不存在的订阅（不应该 panic）
	bus.Unsubscribe("non-existent", 999)
}

func TestEventBus_Unsubscribe_LastSubscriber(t *testing.T) {
	bus := New()

	id, _ := bus.Subscribe("test-topic", 10)

	// 取消最后一个订阅者
	bus.Unsubscribe("test-topic", id)

	// topic 应该被删除
	if _, exists := bus.subs["test-topic"]; exists {
		t.Fatal("topic should be removed when last subscriber unsubscribes")
	}
}

func TestEventBus_ConcurrentAccess(t *testing.T) {
	bus := New()

	var wg sync.WaitGroup
	subscribers := 10
	messages := 100

	// 创建多个订阅者
	channels := make([]<-chan interface{}, subscribers)
	for i := 0; i < subscribers; i++ {
		id, ch := bus.Subscribe("concurrent-topic", 100)
		channels[i] = ch
		_ = id
	}

	// 并发发布消息
	wg.Add(messages)
	for i := 0; i < messages; i++ {
		go func(msg int) {
			defer wg.Done()
			bus.Publish("concurrent-topic", msg)
		}(i)
	}

	// 并发取消订阅
	wg.Add(subscribers / 2)
	for i := 0; i < subscribers/2; i++ {
		go func(id int) {
			defer wg.Done()
			bus.Unsubscribe("concurrent-topic", id+1)
		}(i)
	}

	wg.Wait()

	// 验证剩余订阅者仍能接收消息
	bus.Publish("concurrent-topic", "final-message")

	// 等待一段时间确保消息传递
	time.Sleep(100 * time.Millisecond)
}

func TestEventBus_MultipleTopics(t *testing.T) {
	bus := New()

	_, ch1 := bus.Subscribe("topic1", 10)
	_, ch2 := bus.Subscribe("topic2", 10)

	bus.Publish("topic1", "message1")
	bus.Publish("topic2", "message2")

	select {
	case msg := <-ch1:
		if msg != "message1" {
			t.Errorf("ch1: expected 'message1', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("ch1: message not received")
	}

	select {
	case msg := <-ch2:
		if msg != "message2" {
			t.Errorf("ch2: expected 'message2', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("ch2: message not received")
	}
}

func TestEventBus_SubscribeID_Increment(t *testing.T) {
	bus := New()

	id1, _ := bus.Subscribe("topic", 10)
	id2, _ := bus.Subscribe("topic", 10)
	id3, _ := bus.Subscribe("topic", 10)

	if id2 <= id1 {
		t.Errorf("id2 (%d) should be greater than id1 (%d)", id2, id1)
	}
	if id3 <= id2 {
		t.Errorf("id3 (%d) should be greater than id2 (%d)", id3, id2)
	}
}

func TestEventBus_DroppedCount(t *testing.T) {
	bus := New()

	// 初始计数应该为 0
	if bus.DroppedCount() != 0 {
		t.Errorf("expected initial dropped count 0, got %d", bus.DroppedCount())
	}

	id, ch := bus.Subscribe("test-topic", 1)

	// 填满 channel
	bus.Publish("test-topic", "msg1")
	<-ch // 接收一个，让 channel 有空间

	// 发布多个消息，应该会丢弃一些
	bus.Publish("test-topic", "msg2")
	bus.Publish("test-topic", "msg3")
	bus.Publish("test-topic", "msg4")

	// 验证有消息被丢弃
	dropped := bus.DroppedCount()
	if dropped == 0 {
		t.Error("expected some messages to be dropped, but count is 0")
	}

	// 测试重置
	bus.ResetDroppedCount()
	if bus.DroppedCount() != 0 {
		t.Errorf("expected dropped count 0 after reset, got %d", bus.DroppedCount())
	}

	bus.Unsubscribe("test-topic", id)
}

func TestEventBus_Publish_ConcurrentUnsubscribe(t *testing.T) {
	bus := New()

	// 创建多个订阅者
	ids := make([]int, 10)
	channels := make([]<-chan interface{}, 10)
	for i := 0; i < 10; i++ {
		id, ch := bus.Subscribe("concurrent-topic", 10)
		ids[i] = id
		channels[i] = ch
	}

	// 并发发布和取消订阅
	var wg sync.WaitGroup
	wg.Add(20)

	// 10 个 goroutine 发布消息
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				bus.Publish("concurrent-topic", j)
			}
		}()
	}

	// 10 个 goroutine 取消订阅
	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * 10) // 稍微延迟，确保有发布操作在进行
			bus.Unsubscribe("concurrent-topic", ids[idx])
		}(i)
	}

	wg.Wait()

	// 验证没有 panic，并且可以继续发布
	bus.Publish("concurrent-topic", "final")
	time.Sleep(50 * time.Millisecond)
}
