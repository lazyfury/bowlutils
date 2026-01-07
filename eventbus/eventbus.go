package eventbus

import (
	"sync"
	"sync/atomic"
)

const (
	// DefaultBufferSize is the default buffer size for channels when buffer <= 0
	DefaultBufferSize = 10
)

// EventBus is a thread-safe event bus implementation that allows
// publishers to send events to multiple subscribers.
// When a subscriber's channel buffer is full, messages are dropped
// (non-blocking behavior) to prevent blocking the publisher.
type EventBus struct {
	mu      sync.RWMutex
	subs    map[string]map[int]chan interface{}
	next    int
	dropped int64 // atomic counter for dropped messages
}

// New creates a new EventBus instance.
func New() *EventBus {
	return &EventBus{
		subs: make(map[string]map[int]chan interface{}),
	}
}

// Subscribe subscribes to a topic and returns a subscription ID and a channel.
// If buffer <= 0, DefaultBufferSize will be used.
// The returned channel will be closed when Unsubscribe is called.
func (b *EventBus) Subscribe(topic string, buffer int) (int, <-chan interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if buffer <= 0 {
		buffer = DefaultBufferSize
	}

	ch := make(chan interface{}, buffer)
	if _, ok := b.subs[topic]; !ok {
		b.subs[topic] = make(map[int]chan interface{})
	}

	b.next++
	id := b.next
	b.subs[topic][id] = ch
	return id, ch
}

// Unsubscribe removes a subscription and closes its channel.
// It is safe to call Unsubscribe multiple times with the same id.
func (b *EventBus) Unsubscribe(topic string, id int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if m, ok := b.subs[topic]; ok {
		if ch, ok := m[id]; ok {
			delete(m, id)
			close(ch)
		}
		if len(m) == 0 {
			delete(b.subs, topic)
		}
	}
}

// Publish sends a payload to all subscribers of the given topic.
// If a subscriber's channel buffer is full, the message is dropped
// (non-blocking) to prevent blocking the publisher.
// This method is thread-safe and can be called concurrently.
func (b *EventBus) Publish(topic string, payload interface{}) {
	b.mu.RLock()
	m, ok := b.subs[topic]
	if !ok {
		b.mu.RUnlock()
		return
	}

	// Create a snapshot of channels to avoid holding the lock
	// while sending messages. This prevents potential deadlocks
	// and race conditions when Unsubscribe is called concurrently.
	channels := make([]chan interface{}, 0, len(m))
	for _, ch := range m {
		channels = append(channels, ch)
	}
	b.mu.RUnlock()

	// Send messages outside the lock to minimize lock contention
	for _, ch := range channels {
		select {
		case ch <- payload:
			// Message sent successfully
		default:
			// Channel buffer is full, drop the message
			atomic.AddInt64(&b.dropped, 1)
		}
	}
}

// DroppedCount returns the number of messages that were dropped
// due to full channel buffers since the EventBus was created.
func (b *EventBus) DroppedCount() int64 {
	return atomic.LoadInt64(&b.dropped)
}

// ResetDroppedCount resets the dropped message counter to zero.
func (b *EventBus) ResetDroppedCount() {
	atomic.StoreInt64(&b.dropped, 0)
}
