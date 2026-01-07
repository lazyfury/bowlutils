package files

import (
	"context"
	"io"
	"sync"
)

// 简单的 Processor 注册与工厂
type ProcessorFactory func() Processor

var (
	processorMu sync.RWMutex
	factories   = make(map[string]ProcessorFactory)
)

func RegisterProcessor(name string, f ProcessorFactory) {
	processorMu.Lock()
	defer processorMu.Unlock()
	factories[name] = f
}

func GetProcessor(name string) (Processor, bool) {
	processorMu.RLock()
	defer processorMu.RUnlock()
	f, ok := factories[name]
	if !ok {
		return nil, false
	}
	return f(), true
}

// 一个示例的 NoOpProcessor
type NoOpProcessor struct{}

func (n *NoOpProcessor) Process(ctx context.Context, in io.Reader, meta Metadata, task ProcessTask) (string, Metadata, error) {
	// 不做任何处理，返回空
	return "", Metadata{}, nil
}

func init() {
	RegisterProcessor("noop", func() Processor { return &NoOpProcessor{} })
}
