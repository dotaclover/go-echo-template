package services

import (
	"fmt"
	"myapp/utils"
)

// NotifyService 统一通知服务（聚合多渠道）
type NotifyService struct {
	channels []Notifier
}

func NewNotifyService(channels ...Notifier) *NotifyService {
	return &NotifyService{channels: channels}
}

// AddChannel 动态添加通知渠道
func (s *NotifyService) AddChannel(ch Notifier) {
	s.channels = append(s.channels, ch)
}

// Send 向所有渠道发送通知（失败不中断，记录日志）
func (s *NotifyService) Send(title, content string) {
	for _, ch := range s.channels {
		if err := ch.Send(title, content); err != nil {
			if utils.Logger != nil {
				utils.Logger.Errorf("[notify] %s send failed: %v", ch.Name(), err)
			}
		}
	}
}

// SendTo 向指定渠道发送通知
func (s *NotifyService) SendTo(channelName, title, content string) error {
	for _, ch := range s.channels {
		if ch.Name() == channelName {
			return ch.Send(title, content)
		}
	}
	return fmt.Errorf("channel %s not found", channelName)
}

// Channels 返回已注册的渠道名称列表
func (s *NotifyService) Channels() []string {
	names := make([]string, len(s.channels))
	for i, ch := range s.channels {
		names[i] = ch.Name()
	}
	return names
}
