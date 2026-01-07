package module

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lazyfury/bowlutils/logger"
)

/*
WorkerModule 使用示例:

	// 创建 Worker 模块（4 个 worker）
	workerModule := module.NewWorkerModule(4)

	// 注册到模块管理器
	moduleManager.RegisterModule("WorkerModule", workerModule)

	// 提交任务
	task := module.NewSimpleTask(
		"process-data",
		func(ctx context.Context) error {
			// 执行任务逻辑
			logger.Info("Processing data...")
			return nil
		},
		module.WithPriority(10),           // 设置优先级
		module.WithTimeout(30*time.Second), // 设置超时
		module.WithRetry(3),              // 设置重试次数
	)

	taskID, err := workerModule.SubmitTask(task)
	if err != nil {
		logger.Error("Failed to submit task", "error", err)
	}

	// 查询任务状态
	taskInfo, exists := workerModule.GetTaskInfo(taskID)
	if exists {
		logger.Info("Task status", "status", taskInfo.Status)
	}

	// 取消任务（仅对未开始的任务有效）
	workerModule.CancelTask(taskID)
*/

// Task 接口定义任务需要实现的方法
type Task interface {
	// Execute 执行任务，返回错误
	Execute(ctx context.Context) error
	// Name 返回任务名称
	Name() string
	// Priority 返回任务优先级，数字越大优先级越高
	Priority() int
	// Timeout 返回任务超时时间，0表示不超时
	Timeout() time.Duration
	// Retry 返回任务重试次数
	Retry() int
}

// TaskStatus 任务状态
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusRunning
	TaskStatusCompleted
	TaskStatusFailed
	TaskStatusCancelled
)

const (
	TaskStatusPendingStr   = "pending"
	TaskStatusRunningStr   = "running"
	TaskStatusCompletedStr = "completed"
	TaskStatusFailedStr    = "failed"
	TaskStatusCancelledStr = "cancelled"
)

var TaskStatusStrMap = map[TaskStatus]string{
	TaskStatusPending:   TaskStatusPendingStr,
	TaskStatusRunning:   TaskStatusRunningStr,
	TaskStatusCompleted: TaskStatusCompletedStr,
	TaskStatusFailed:    TaskStatusFailedStr,
	TaskStatusCancelled: TaskStatusCancelledStr,
}

// TaskInfo 任务信息
type TaskInfo struct {
	ID        string
	Task      Task
	Status    TaskStatus
	StatusStr string
	CreatedAt time.Time
	StartedAt time.Time
	EndedAt   time.Time
	Error     error
	Retries   int
}

// WorkerModule Worker Pool 模块，用于并发执行任务
type WorkerModule struct {
	workerCount int
	submitQueue chan *TaskInfo // 任务提交队列
	taskQueue   chan *TaskInfo // Worker 消费队列
	quit        chan bool
	tasks       map[string]*TaskInfo
	tasksMutex  sync.RWMutex
	wg          sync.WaitGroup
}

// NewWorkerModule 创建新的 Worker 模块
// workerCount: worker 数量，建议设置为 CPU 核心数或稍大
func NewWorkerModule(workerCount int) *WorkerModule {
	if workerCount <= 0 {
		workerCount = 1
	}
	return &WorkerModule{
		workerCount: workerCount,
		submitQueue: make(chan *TaskInfo, 100), // 任务提交队列缓冲区
		taskQueue:   make(chan *TaskInfo, 100), // Worker 消费队列缓冲区
		quit:        make(chan bool),
		tasks:       make(map[string]*TaskInfo),
	}
}

// SubmitTask 提交任务到队列
func (wm *WorkerModule) SubmitTask(task Task) (string, error) {
	taskID := generateTaskID()
	taskInfo := &TaskInfo{
		ID:        taskID,
		Task:      task,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		Retries:   0,
	}

	wm.tasksMutex.Lock()
	wm.tasks[taskID] = taskInfo
	wm.tasksMutex.Unlock()

	select {
	case wm.submitQueue <- taskInfo:
		logger.Info("Task submitted", "[task_id]", taskID, "[task_name]", task.Name(), "[priority]", task.Priority())
		return taskID, nil
	default:
		return "", fmt.Errorf("task queue is full")
	}
}

// GetTaskInfo 获取任务信息
func (wm *WorkerModule) GetTaskInfo(taskID string) (*TaskInfo, bool) {
	wm.tasksMutex.RLock()
	defer wm.tasksMutex.RUnlock()
	if wm.tasks == nil {
		return nil, false
	}
	taskInfo, exists := wm.tasks[taskID]
	if !exists {
		return nil, false
	}
	taskInfo.StatusStr = TaskStatusStrMap[taskInfo.Status]
	return taskInfo, true
}

// CancelTask 取消任务（仅对未开始的任务有效）
func (wm *WorkerModule) CancelTask(taskID string) bool {
	wm.tasksMutex.Lock()
	defer wm.tasksMutex.Unlock()
	taskInfo, exists := wm.tasks[taskID]
	if !exists || taskInfo.Status != TaskStatusPending {
		return false
	}
	taskInfo.Status = TaskStatusCancelled
	taskInfo.EndedAt = time.Now()
	return true
}

// Start 启动 Worker 模块
func (wm *WorkerModule) Start(ctx context.Context, wg *sync.WaitGroup) error {
	// 启动 worker goroutines
	for i := 0; i < wm.workerCount; i++ {
		wg.Add(1)
		go wm.worker(ctx, i, wg)
	}

	// 启动任务调度器（按优先级排序）
	wg.Add(1)
	go wm.scheduler(ctx, wg)

	logger.Info("WorkerModule started", "worker_count", wm.workerCount)
	return nil
}

// Stop 停止 Worker 模块
func (wm *WorkerModule) Stop() error {
	logger.Info("WorkerModule stopping")
	close(wm.quit)
	close(wm.submitQueue)
	close(wm.taskQueue)
	wm.wg.Wait()
	logger.Info("WorkerModule stopped")
	return nil
}

// worker 工作协程
func (wm *WorkerModule) worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("Worker started", "worker_id", id)

	for {
		select {
		case taskInfo, ok := <-wm.taskQueue:
			if !ok {
				logger.Info("Worker stopped", "worker_id", id)
				return
			}

			// 检查任务是否已取消
			if taskInfo.Status == TaskStatusCancelled {
				continue
			}

			// 执行任务
			wm.executeTask(ctx, taskInfo)

		case <-wm.quit:
			logger.Info("Worker quitting", "worker_id", id)
			return
		case <-ctx.Done():
			logger.Info("Worker context cancelled", "worker_id", id)
			return
		}
	}
}

// scheduler 任务调度器，按优先级排序任务
func (wm *WorkerModule) scheduler(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// 优先级队列（简单的实现，可以后续优化为堆）
	pendingTasks := make([]*TaskInfo, 0)
	ticker := time.NewTicker(100 * time.Millisecond) // 定期处理待处理任务
	defer ticker.Stop()

	for {
		select {
		case taskInfo, ok := <-wm.submitQueue:
			if !ok {
				// 队列关闭，处理剩余任务
				wm.flushPendingTasks(pendingTasks)
				return
			}

			// 添加到待处理队列
			pendingTasks = append(pendingTasks, taskInfo)

		case <-ticker.C:
			// 定期按优先级发送任务
			if len(pendingTasks) > 0 {
				wm.flushPendingTasks(pendingTasks)
				pendingTasks = pendingTasks[:0]
			}

		case <-wm.quit:
			wm.flushPendingTasks(pendingTasks)
			return
		case <-ctx.Done():
			wm.flushPendingTasks(pendingTasks)
			return
		}
	}
}

// flushPendingTasks 按优先级刷新待处理任务到 worker 队列
func (wm *WorkerModule) flushPendingTasks(tasks []*TaskInfo) {
	if len(tasks) == 0 {
		return
	}

	// 简单的优先级排序（可以优化为堆）
	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].Task.Priority() < tasks[j].Task.Priority() {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}

	// 发送到 worker 队列
	for _, task := range tasks {
		select {
		case wm.taskQueue <- task:
		default:
			logger.Warn("Task queue full, dropping task", "task_id", task.ID)
		}
	}
}

// executeTask 执行任务
func (wm *WorkerModule) executeTask(ctx context.Context, taskInfo *TaskInfo) {
	wm.tasksMutex.Lock()
	taskInfo.Status = TaskStatusRunning
	taskInfo.StartedAt = time.Now()
	wm.tasksMutex.Unlock()

	// 设置任务超时
	taskCtx := ctx
	if timeout := taskInfo.Task.Timeout(); timeout > 0 {
		var cancel context.CancelFunc
		taskCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// 执行任务
	err := taskInfo.Task.Execute(taskCtx)

	wm.tasksMutex.Lock()
	taskInfo.EndedAt = time.Now()

	if err != nil {
		// 检查是否需要重试
		if taskInfo.Retries < taskInfo.Task.Retry() {
			taskInfo.Retries++
			taskInfo.Status = TaskStatusPending
			logger.Warn("Task failed, retrying", "task_id", taskInfo.ID, "retries", taskInfo.Retries, "error", err.Error())

			// 重新加入队列
			select {
			case wm.submitQueue <- taskInfo:
			default:
				taskInfo.Status = TaskStatusFailed
				taskInfo.Error = err
				logger.Error("Task retry failed, queue full", "[task_id]", taskInfo.ID)
			}
		} else {
			taskInfo.Status = TaskStatusFailed
			taskInfo.Error = err
			logger.Error("Task execution failed", "[task_id]", taskInfo.ID, "[error]", err.Error())
		}
	} else {
		taskInfo.Status = TaskStatusCompleted
		logger.Info("Task completed ", "[task_id]", taskInfo.ID, "[duration]", taskInfo.EndedAt.Sub(taskInfo.StartedAt))
	}
	wm.tasksMutex.Unlock()
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d_%d", time.Now().UnixNano(), time.Now().Unix())
}

// SimpleTask 简单的任务实现示例
type SimpleTask struct {
	name     string
	priority int
	timeout  time.Duration
	retry    int
	handler  func(ctx context.Context) error
}

// NewSimpleTask 创建简单任务
func NewSimpleTask(name string, handler func(ctx context.Context) error, opts ...TaskOption) *SimpleTask {
	task := &SimpleTask{
		name:     name,
		priority: 0,
		timeout:  0,
		retry:    0,
		handler:  handler,
	}

	for _, opt := range opts {
		opt(task)
	}

	return task
}

// TaskOption 任务选项
type TaskOption func(*SimpleTask)

// WithPriority 设置任务优先级
func WithPriority(priority int) TaskOption {
	return func(t *SimpleTask) {
		t.priority = priority
	}
}

// WithTimeout 设置任务超时
func WithTimeout(timeout time.Duration) TaskOption {
	return func(t *SimpleTask) {
		t.timeout = timeout
	}
}

// WithRetry 设置重试次数
func WithRetry(retry int) TaskOption {
	return func(t *SimpleTask) {
		t.retry = retry
	}
}

func (st *SimpleTask) Execute(ctx context.Context) error {
	return st.handler(ctx)
}

func (st *SimpleTask) Name() string {
	return st.name
}

func (st *SimpleTask) Priority() int {
	return st.priority
}

func (st *SimpleTask) Timeout() time.Duration {
	return st.timeout
}

func (st *SimpleTask) Retry() int {
	return st.retry
}
