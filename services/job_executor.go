package services

import (
	"context"
	"fmt"
	"myapp/utils"
	"sync"
	"time"
)

// JobExecutor 通用任务执行器（Worker Pool）
type JobExecutor struct {
	queue   QueueInterface
	lock    LockInterface
	workers int
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewJobExecutor(queue QueueInterface, lock LockInterface, workers int) *JobExecutor {
	if workers <= 0 {
		workers = 1
	}
	return &JobExecutor{
		queue:   queue,
		lock:    lock,
		workers: workers,
		stopCh:  make(chan struct{}),
	}
}

// RegisterHandler 注册任务处理器
func (e *JobExecutor) RegisterHandler(handler TaskHandler) {
	e.queue.RegisterHandler(handler)
	utils.Logger.Infof("JobExecutor: registered handler [%s]", handler.GetTaskType())
}

// Start 启动 worker pool
func (e *JobExecutor) Start() {
	utils.Logger.Infof("JobExecutor starting (%d workers)", e.workers)
	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go e.worker(i)
	}
}

// Stop 优雅停止
func (e *JobExecutor) Stop() {
	utils.Logger.Info("JobExecutor stopping...")
	close(e.stopCh)
	e.wg.Wait()
	utils.Logger.Info("JobExecutor stopped")
}

func (e *JobExecutor) worker(id int) {
	defer e.wg.Done()
	for {
		select {
		case <-e.stopCh:
			return
		default:
			e.processOne(id)
		}
	}
}

func (e *JobExecutor) processOne(workerID int) {
	task, err := e.queue.PopImmediate(QueueImmediate, QueueProcessing, 5)
	if err != nil {
		utils.Logger.Errorf("worker %d: pop error: %v", workerID, err)
		return
	}
	if task == nil {
		return // timeout, no task
	}

	// 获取任务锁（防重复执行）
	lockKey := "job:lock:" + task.GetID()
	lock, err := e.lock.TryObtain(lockKey, 10*time.Minute)
	if err != nil {
		utils.Logger.Warnf("worker %d: task %s locked, skip", workerID, task.GetID())
		return
	}
	defer lock.Release()

	// 查找 handler
	handler, err := e.queue.GetHandler(task.GetType())
	if err != nil {
		utils.Logger.Errorf("worker %d: no handler for %s", workerID, task.GetType())
		e.queue.CompleteTask(QueueProcessing, task)
		return
	}

	// 执行
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Now()
	if err := handler.Handle(ctx, task); err != nil {
		utils.Logger.Errorf("worker %d: task %s failed (%v): %v", workerID, task.GetID(), time.Since(start), err)
	} else {
		utils.Logger.Infof("worker %d: task %s done (%v)", workerID, task.GetID(), time.Since(start))
	}

	e.queue.CompleteTask(QueueProcessing, task)
}

// Stats 获取队列统计
func (e *JobExecutor) Stats() map[string]interface{} {
	immediate, _ := e.queue.GetQueueSize(QueueImmediate)
	scheduled, _ := e.queue.GetQueueSize(QueueScheduled)
	retry, _ := e.queue.GetQueueSize(QueueRetry)
	processing, _ := e.queue.GetQueueSize(QueueProcessing)

	return map[string]interface{}{
		"workers": e.workers,
		"queues": map[string]int64{
			"immediate":  immediate,
			"scheduled":  scheduled,
			"retry":      retry,
			"processing": processing,
		},
	}
}

// ============================================================================
// 示例 Handler（可删除，仅供参考）
// ============================================================================

// ExampleHandler 示例任务处理器
type ExampleHandler struct{}

func (h *ExampleHandler) GetTaskType() string              { return "example" }
func (h *ExampleHandler) CanHandle(taskType string) bool   { return taskType == "example" }
func (h *ExampleHandler) Handle(ctx context.Context, task Task) error {
	fmt.Printf("[ExampleHandler] processing task %s, payload: %s\n", task.GetID(), string(task.GetPayload()))
	return nil
}
