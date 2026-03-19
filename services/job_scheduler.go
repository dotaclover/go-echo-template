package services

import (
	"myapp/utils"
	"sync"
	"time"
)

// JobScheduler 定时任务调度器
// 周期性扫描 scheduled/retry 队列，将到期任务移入 immediate 队列
type JobScheduler struct {
	queue    QueueInterface
	interval time.Duration
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

func NewJobScheduler(queue QueueInterface, interval time.Duration) *JobScheduler {
	if interval == 0 {
		interval = 10 * time.Second
	}
	return &JobScheduler{queue: queue, interval: interval, stopCh: make(chan struct{})}
}

// Start 启动调度器
func (s *JobScheduler) Start() {
	utils.Logger.Infof("JobScheduler starting (interval: %v)", s.interval)
	s.wg.Add(1)
	go s.run()
}

// Stop 停止调度器
func (s *JobScheduler) Stop() {
	utils.Logger.Info("JobScheduler stopping...")
	close(s.stopCh)
	s.wg.Wait()
	utils.Logger.Info("JobScheduler stopped")
}

func (s *JobScheduler) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.scan() // 启动时先扫一次

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.scan()
		}
	}
}

func (s *JobScheduler) scan() {
	s.scanQueue(QueueScheduled)
	s.scanQueue(QueueRetry)
}

func (s *JobScheduler) scanQueue(queueName string) {
	tasks, err := s.queue.ScanScheduled(queueName)
	if err != nil {
		utils.Logger.Errorf("JobScheduler: scan %s error: %v", queueName, err)
		return
	}
	if len(tasks) == 0 {
		return
	}

	utils.Logger.Infof("JobScheduler: found %d tasks in %s", len(tasks), queueName)

	for _, task := range tasks {
		if err := s.queue.PushImmediate(QueueImmediate, task); err != nil {
			utils.Logger.Errorf("JobScheduler: push task %s failed: %v", task.GetID(), err)
			continue
		}
		if err := s.queue.RemoveScheduled(queueName, task); err != nil {
			utils.Logger.Errorf("JobScheduler: remove task %s from %s failed: %v", task.GetID(), queueName, err)
		}
	}
}

// Stats 获取调度器统计
func (s *JobScheduler) Stats() map[string]interface{} {
	immediate, _ := s.queue.GetQueueSize(QueueImmediate)
	scheduled, _ := s.queue.GetQueueSize(QueueScheduled)
	retry, _ := s.queue.GetQueueSize(QueueRetry)
	processing, _ := s.queue.GetQueueSize(QueueProcessing)

	return map[string]interface{}{
		"interval": s.interval.String(),
		"queues": map[string]int64{
			"immediate":  immediate,
			"scheduled":  scheduled,
			"retry":      retry,
			"processing": processing,
		},
	}
}
