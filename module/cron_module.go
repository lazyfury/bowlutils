package module

import (
	"context"
	"sync"
	"time"

	"github.com/lazyfury/bowlutils/logger"

	"github.com/robfig/cron/v3"
)

// Job Cron 任务定义
type Job struct {
	spec string
	job  func()
}

// CornModule 基于 cron 的定时任务模块
type CornModule struct {
	cron *cron.Cron
	quit chan bool
	Jobs []Job
}

// NewCornModule 创建新的 Cron 模块
func NewCornModule(c *cron.Cron, jobs ...Job) *CornModule {
	return &CornModule{
		cron: c,
		quit: make(chan bool),
		Jobs: jobs,
	}
}

// AddJob 添加 Cron 任务
func (cm *CornModule) AddJob(spec string, job func()) {
	cm.Jobs = append(cm.Jobs, Job{spec: spec, job: job})
}

// Start 启动 Cron 模块
func (cm *CornModule) Start(ctx context.Context, wg *sync.WaitGroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		cm.cron.AddFunc("@every 1m", func() {
			logger.Info("CornModule cron job executed", "time", time.Now().Format("2006-01-02 15:04:05"))
		})
		for _, job := range cm.Jobs {
			cm.cron.AddFunc(job.spec, job.job)
		}
		cm.cron.Start()
		<-cm.quit
	}()
	return nil
}

// Stop 停止 Cron 模块
func (cm *CornModule) Stop() error {
	logger.Info("CornModule stopping")
	cm.cron.Stop()
	close(cm.quit)
	return nil
}
