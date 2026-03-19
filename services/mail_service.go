package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// MailService SMTP 邮件服务
type MailService struct {
	host     string
	port     int
	username string
	password string
	from     string
	useTLS   bool
}

// MailConfig 邮件配置
type MailConfig struct {
	Host     string // SMTP 服务器
	Port     int    // 端口（465=TLS, 587=STARTTLS, 25=明文）
	Username string
	Password string
	From     string // 发件人地址
	UseTLS   bool   // 是否使用 TLS
}

func NewMailService(cfg MailConfig) *MailService {
	return &MailService{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
		useTLS:   cfg.UseTLS,
	}
}

// Send 发送纯文本邮件
func (s *MailService) Send(to []string, subject, body string) error {
	return s.SendHTML(to, subject, body, false)
}

// SendHTML 发送邮件（支持 HTML）
func (s *MailService) SendHTML(to []string, subject, body string, isHTML bool) error {
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + strings.Join(to, ","),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: " + contentType + "; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	if s.useTLS {
		return s.sendWithTLS(addr, auth, to, []byte(msg))
	}
	return smtp.SendMail(addr, auth, s.from, to, []byte(msg))
}

func (s *MailService) sendWithTLS(addr string, auth smtp.Auth, to []string, msg []byte) error {
	tlsConfig := &tls.Config{ServerName: s.host}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("SMTP client failed: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	if err = client.Mail(s.from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}
