package services

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ============================================================================
// SQLite 队列表模型
// ============================================================================

// QueueTask SQLite 队列任务表
type QueueTask struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	QueueName   string     `gorm:"size:100;not null;index:idx_queue_status" json:"queue_name"`
	TaskID      string     `gorm:"size:200;not null;index" json:"task_id"`
	TaskType    string     `gorm:"size:100;not null" json:"task_type"`
	Priority    int        `gorm:"default:0" json:"priority"`
	Payload     string     `gorm:"type:text" json:"payload"`
	Status      int        `gorm:"default:0;index:idx_queue_status" json:"status"`
	ScheduledAt *time.Time `gorm:"index" json:"scheduled_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (QueueTask) TableName() string { return "queue_tasks" }

const (
	taskStatusPending    = 0
	taskStatusProcessing = 1
)

// ============================================================================
// SQLite 队列服务实现
// ============================================================================

// SQLiteQueueService SQLite 实现的通用队列
type SQLiteQueueService struct {
	db       *gorm.DB
	handlers map[string]TaskHandler
}

func NewSQLiteQueueService(db *gorm.DB) *SQLiteQueueService {
	return &SQLiteQueueService{db: db, handlers: make(map[string]TaskHandler)}
}

func (s *SQLiteQueueService) RegisterHandler(handler TaskHandler) {
	s.handlers[handler.GetTaskType()] = handler
}

func (s *SQLiteQueueService) GetHandler(taskType string) (TaskHandler, error) {
	h, ok := s.handlers[taskType]
	if !ok {
		return nil, fmt.Errorf("no handler for task type: %s", taskType)
	}
	return h, nil
}

// PushImmediate 推送立即任务
func (s *SQLiteQueueService) PushImmediate(queueName string, task Task) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}
	return s.db.Create(&QueueTask{
		QueueName: queueName,
		TaskID:    task.GetID(),
		TaskType:  task.GetType(),
		Priority:  task.GetPriority(),
		Payload:   string(payload),
		Status:    taskStatusPending,
		CreatedAt: time.Now(),
	}).Error
}

// PopImmediate 弹出任务（阻塞等待，timeout 秒）
func (s *SQLiteQueueService) PopImmediate(queueName, processingQueue string, timeout int) (Task, error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for time.Now().Before(deadline) {
		var qt QueueTask
		err := s.db.Transaction(func(tx *gorm.DB) error {
			result := tx.Where("queue_name = ? AND status = ?", queueName, taskStatusPending).
				Order("priority DESC, id ASC").First(&qt)
			if result.Error != nil {
				return result.Error
			}
			return tx.Model(&qt).Updates(map[string]interface{}{
				"status":     taskStatusProcessing,
				"queue_name": processingQueue,
			}).Error
		})

		if err == gorm.ErrRecordNotFound {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err != nil {
			return nil, err
		}

		var task BaseTask
		if err := json.Unmarshal([]byte(qt.Payload), &task); err != nil {
			s.db.Delete(&qt)
			continue
		}
		return &task, nil
	}
	return nil, nil // timeout
}

// PushScheduled 推送定时任务
func (s *SQLiteQueueService) PushScheduled(queueName string, task Task, scheduledAtUnix int64) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}
	t := time.Unix(scheduledAtUnix, 0)
	return s.db.Create(&QueueTask{
		QueueName:   queueName,
		TaskID:      task.GetID(),
		TaskType:    task.GetType(),
		Priority:    task.GetPriority(),
		Payload:     string(payload),
		Status:      taskStatusPending,
		ScheduledAt: &t,
		CreatedAt:   time.Now(),
	}).Error
}

// ScanScheduled 扫描到期的定时任务
func (s *SQLiteQueueService) ScanScheduled(queueName string) ([]Task, error) {
	now := time.Now()
	var items []QueueTask
	err := s.db.Where("queue_name = ? AND status = ? AND scheduled_at <= ?",
		queueName, taskStatusPending, now).Find(&items).Error
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, len(items))
	for _, qt := range items {
		var t BaseTask
		if json.Unmarshal([]byte(qt.Payload), &t) == nil {
			tasks = append(tasks, &t)
		}
	}
	return tasks, nil
}

// RemoveScheduled 移除定时任务
func (s *SQLiteQueueService) RemoveScheduled(queueName string, task Task) error {
	return s.db.Where("queue_name = ? AND task_id = ? AND status = ?",
		queueName, task.GetID(), taskStatusPending).Delete(&QueueTask{}).Error
}

// CompleteTask 完成任务（从 processing 队列删除）
func (s *SQLiteQueueService) CompleteTask(processingQueue string, task Task) error {
	return s.db.Where("queue_name = ? AND task_id = ?",
		processingQueue, task.GetID()).Delete(&QueueTask{}).Error
}

// GetQueueSize 获取队列待处理任务数
func (s *SQLiteQueueService) GetQueueSize(queueName string) (int64, error) {
	var count int64
	err := s.db.Model(&QueueTask{}).
		Where("queue_name = ? AND status = ?", queueName, taskStatusPending).
		Count(&count).Error
	return count, err
}

// ClearQueue 清空队列
func (s *SQLiteQueueService) ClearQueue(queueName string) error {
	return s.db.Where("queue_name = ?", queueName).Delete(&QueueTask{}).Error
}
