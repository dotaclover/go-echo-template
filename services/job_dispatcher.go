package services

import (
	"encoding/json"
	"fmt"
	"time"
)

// JobDispatcher 任务分发器
// 业务代码通过 Dispatcher 将任务投递到队列
type JobDispatcher struct {
	queue QueueInterface
}

func NewJobDispatcher(queue QueueInterface) *JobDispatcher {
	return &JobDispatcher{queue: queue}
}

// Dispatch 立即分发任务
func (d *JobDispatcher) Dispatch(taskType string, payload interface{}) error {
	return d.DispatchWithPriority(taskType, payload, 0)
}

// DispatchWithPriority 带优先级分发
func (d *JobDispatcher) DispatchWithPriority(taskType string, payload interface{}, priority int) error {
	task, err := d.buildTask(taskType, payload, priority)
	if err != nil {
		return err
	}
	return d.queue.PushImmediate(QueueImmediate, task)
}

// DispatchLater 延迟分发
func (d *JobDispatcher) DispatchLater(taskType string, payload interface{}, delay time.Duration) error {
	return d.DispatchAt(taskType, payload, time.Now().Add(delay))
}

// DispatchAt 定时分发
func (d *JobDispatcher) DispatchAt(taskType string, payload interface{}, scheduledAt time.Time) error {
	task, err := d.buildTask(taskType, payload, 0)
	if err != nil {
		return err
	}
	return d.queue.PushScheduled(QueueScheduled, task, scheduledAt.Unix())
}

func (d *JobDispatcher) buildTask(taskType string, payload interface{}, priority int) (*BaseTask, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload failed: %w", err)
	}
	return &BaseTask{
		ID:       fmt.Sprintf("%s-%d", taskType, time.Now().UnixNano()),
		Type:     taskType,
		Priority: priority,
		Payload:  data,
	}, nil
}

// Stats 队列统计
func (d *JobDispatcher) Stats() map[string]int64 {
	immediate, _ := d.queue.GetQueueSize(QueueImmediate)
	scheduled, _ := d.queue.GetQueueSize(QueueScheduled)
	retry, _ := d.queue.GetQueueSize(QueueRetry)
	processing, _ := d.queue.GetQueueSize(QueueProcessing)

	return map[string]int64{
		"immediate":  immediate,
		"scheduled":  scheduled,
		"retry":      retry,
		"processing": processing,
	}
}
