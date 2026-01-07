package module

import (
	"context"
	"sync"
	"time"

	"github.com/lazyfury/bowlutils/logger"
)

// TickModule 基于 time.Ticker 的定时任务模块
type TickModule struct {
	ticker *time.Ticker
	quit   chan bool
}

// NewTickModule 创建新的 Tick 模块
func NewTickModule(interval time.Duration) *TickModule {
	return &TickModule{
		ticker: time.NewTicker(interval),
		quit:   make(chan bool),
	}
}

// Start 启动 Tick 模块
func (tm *TickModule) Start(ctx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-tm.ticker.C:
				// 执行定时任务
				logger.Info("TickModule tick executed", "time", time.Now().String())
			case <-tm.quit:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

// Stop 停止 Tick 模块
func (tm *TickModule) Stop() error {
	close(tm.quit)
	tm.ticker.Stop()
	return nil
}
