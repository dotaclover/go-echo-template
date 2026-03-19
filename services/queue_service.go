package services

import "context"

// ============================================================================
// 队列任务接口
// ============================================================================

// Task 通用任务接口
type Task interface {
	GetID() string
	GetType() string
	GetPriority() int
	GetPayload() []byte
	SetPayload([]byte) error
}

// BaseTask 基础任务实现
type BaseTask struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Priority int    `json:"priority"`
	Payload  []byte `json:"payload"`
}

func (t *BaseTask) GetID() string       { return t.ID }
func (t *BaseTask) GetType() string     { return t.Type }
func (t *BaseTask) GetPriority() int    { return t.Priority }
func (t *BaseTask) GetPayload() []byte  { return t.Payload }
func (t *BaseTask) SetPayload(d []byte) error { t.Payload = d; return nil }

// TaskHandler 任务处理器接口
type TaskHandler interface {
	Handle(ctx context.Context, task Task) error
	CanHandle(taskType string) bool
	GetTaskType() string
}

// ============================================================================
// 通用队列接口
// ============================================================================

// QueueInterface 通用队列服务接口
// SQLite 和 Redis 均可实现此接口
type QueueInterface interface {
	// Handler 注册
	RegisterHandler(handler TaskHandler)
	GetHandler(taskType string) (TaskHandler, error)

	// 立即任务
	PushImmediate(queueName string, task Task) error
	PopImmediate(queueName, processingQueue string, timeout int) (Task, error)

	// 定时任务
	PushScheduled(queueName string, task Task, scheduledAtUnix int64) error
	ScanScheduled(queueName string) ([]Task, error)
	RemoveScheduled(queueName string, task Task) error

	// 完成 / 统计
	CompleteTask(processingQueue string, task Task) error
	GetQueueSize(queueName string) (int64, error)
	ClearQueue(queueName string) error
}

// 队列名称常量
const (
	QueueImmediate  = "queue:immediate"
	QueueScheduled  = "queue:scheduled"
	QueueRetry      = "queue:retry"
	QueueProcessing = "queue:processing"
)
