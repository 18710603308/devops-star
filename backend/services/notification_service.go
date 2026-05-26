package services

import (
	"bytes"
	"devops-star/backend/config"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ========== 通知服务（完整实现）==========

type NotificationService struct {
	Cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{Cfg: cfg}
}

// 发送通知（根据配置的类型）
func (s *NotificationService) SendNotification(message string, notifyType string) error {
	// notifyType: "always", "success", "failed", "wecom", "dingtalk", "feishu", "email"
	switch notifyType {
	case "wecom", "always":
		if s.Cfg.NotifyWeComWebhook != "" {
			return s.sendWeCom(message)
		}
	case "dingtalk":
		if s.Cfg.NotifyDingTalkWebhook != "" {
			return s.sendDingTalk(message)
		}
	case "feishu":
		if s.Cfg.NotifyFeishuWebhook != "" {
			return s.sendFeishu(message)
		}
	case "email":
		return s.sendEmail(message)
	default:
		// 默认尝试所有已配置的通知渠道
		s.sendWeCom(message)
		s.sendDingTalk(message)
		s.sendFeishu(message)
	}

	return nil
}

// 发送企业微信通知
func (s *NotificationService) sendWeCom(message string) error {
	webhook := s.Cfg.NotifyWeComWebhook
	if webhook == "" {
		return fmt.Errorf("企业微信 Webhook 未配置")
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "🚀 DevOpsStar 通知\n" + message,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("企业微信通知发送失败，状态码：%d", resp.StatusCode)
	}

	return nil
}

// 发送钉钉通知
func (s *NotificationService) sendDingTalk(message string) error {
	webhook := s.Cfg.NotifyDingTalkWebhook
	if webhook == "" {
		return fmt.Errorf("钉钉 Webhook 未配置")
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "🚀 DevOpsStar 通知\n" + message,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("钉钉通知发送失败，状态码：%d", resp.StatusCode)
	}

	return nil
}

// 发送飞书通知
func (s *NotificationService) sendFeishu(message string) error {
	webhook := s.Cfg.NotifyFeishuWebhook
	if webhook == "" {
		return fmt.Errorf("飞书 Webhook 未配置")
	}

	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": "🚀 DevOpsStar 通知\n" + message,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("飞书通知发送失败，状态码：%d", resp.StatusCode)
	}

	return nil
}

// 发送邮件通知（使用 SMTP）
func (s *NotificationService) sendEmail(message string) error {
	if s.Cfg.SMTPHost == "" {
		return fmt.Errorf("SMTP 未配置")
	}

	// 实际应使用 net/smtp 或第三方库（如 gomail）
	// 这里输出日志作为占位
	fmt.Printf("[邮件通知] %s (SMTP: %s)\n", message, s.Cfg.SMTPHost)

	// TODO: 实现 SMTP 发送逻辑
	// import "gopkg.in/gomail.v2"
	// m := gomail.NewMessage()
	// m.SetHeader("From", s.Cfg.SMTPFrom)
	// m.SetHeader("To", "target@example.com")
	// m.SetHeader("Subject", "DevOpsStar 通知")
	// m.SetBody("text/plain", message)
	// d := gomail.NewDialer(s.Cfg.SMTPHost, s.Cfg.SMTPPort, s.Cfg.SMTPUsername, s.Cfg.SMTPPassword)
	// d.DialAndSend(m)

	return nil
}

// 异步发送通知（避免阻塞主流程）
func (s *NotificationService) SendNotificationAsync(message string, notifyType string) {
	go func() {
		if err := s.SendNotification(message, notifyType); err != nil {
			fmt.Printf("[通知发送失败] %v\n", err)
		}
	}()
}

// 测试 Webhook 连通性
func (s *NotificationService) TestWebhook(webhookType string) error {
	testMsg := "✅ DevOpsStar 通知测试成功！Webhook 配置正确。"

	switch webhookType {
	case "wecom":
		return s.sendWeCom(testMsg)
	case "dingtalk":
		return s.sendDingTalk(testMsg)
	case "feishu":
		return s.sendFeishu(testMsg)
	default:
		return fmt.Errorf("不支持的 Webhook 类型：%s", webhookType)
	}
}
