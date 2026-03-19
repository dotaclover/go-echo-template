package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FeishuNotifier 飞书 Webhook 通知
type FeishuNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewFeishuNotifier(webhookURL string) *FeishuNotifier {
	return &FeishuNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (f *FeishuNotifier) Name() string { return "feishu" }

func (f *FeishuNotifier) Send(title, content string) error {
	payload := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": title,
				},
			},
			"elements": []map[string]interface{}{
				{
					"tag": "markdown",
					"content": content,
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := f.client.Post(f.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("feishu send failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("feishu error [%d]: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
