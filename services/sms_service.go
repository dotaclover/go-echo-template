package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SMSService 短信发送服务（预留接口 + 阿里云实现示例）
type SMSService struct {
	provider SMSSender
}

func NewSMSService(provider SMSSender) *SMSService {
	return &SMSService{provider: provider}
}

func (s *SMSService) Send(phone, templateID string, params map[string]string) error {
	return s.provider.Send(phone, templateID, params)
}

// ============================================================================
// 阿里云短信实现（示例，需按实际签名方式调整）
// ============================================================================

type AliyunSMS struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	Endpoint        string
}

func NewAliyunSMS(keyID, keySecret, signName string) *AliyunSMS {
	return &AliyunSMS{
		AccessKeyID:     keyID,
		AccessKeySecret: keySecret,
		SignName:        signName,
		Endpoint:        "https://dysmsapi.aliyuncs.com",
	}
}

func (a *AliyunSMS) Send(phone, templateID string, params map[string]string) error {
	// 简化示例：实际使用时需要按阿里云 API 签名规范处理
	paramsJSON, _ := json.Marshal(params)
	body := map[string]string{
		"PhoneNumbers":  phone,
		"SignName":      a.SignName,
		"TemplateCode":  templateID,
		"TemplateParam": string(paramsJSON),
	}
	bodyJSON, _ := json.Marshal(body)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(a.Endpoint, "application/json", bytes.NewReader(bodyJSON))
	if err != nil {
		return fmt.Errorf("sms request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sms failed [%d]: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// ============================================================================
// Mock 短信（开发测试用）
// ============================================================================

type MockSMS struct{}

func NewMockSMS() *MockSMS { return &MockSMS{} }

func (m *MockSMS) Send(phone, templateID string, params map[string]string) error {
	fmt.Printf("[MockSMS] phone=%s template=%s params=%v\n", phone, templateID, params)
	return nil
}
