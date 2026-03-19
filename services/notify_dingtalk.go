package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DingtalkNotifier 钉钉 Webhook 通知
type DingtalkNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewDingtalkNotifier(webhookURL string) *DingtalkNotifier {
	return &DingtalkNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (d *DingtalkNotifier) Name() string { return "dingtalk" }

func (d *DingtalkNotifier) Send(title, content string) error {
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  fmt.Sprintf("### %s\n%s", title, content),
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("dingtalk send failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("dingtalk error [%d]: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
