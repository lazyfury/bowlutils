package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"
)

// Config Email 配置
type Config struct {
	Host     string // SMTP 服务器地址，如 smtp.gmail.com
	Port     int    // SMTP 端口，如 587
	Username string // 发送邮箱用户名
	Password string // 发送邮箱密码或授权码
	From     string // 发送者邮箱地址
	FromName string // 发送者名称（可选）
	TLS      bool   // 是否使用 TLS
}

// Message Email 消息
type Message struct {
	To      []string // 收件人列表
	Cc      []string // 抄送列表（可选）
	Bcc     []string // 密送列表（可选）
	Subject string   // 主题
	Body    string   // 正文（纯文本）
	HTML    string   // 正文（HTML，如果提供则优先使用）
}

// Sender Email 发送器接口
type Sender interface {
	// Send 同步发送邮件
	Send(ctx context.Context, msg *Message) error
}

// SMTPSender SMTP 发送器实现
type SMTPSender struct {
	config *Config
}

var _ Sender = (*SMTPSender)(nil)

// NewSMTPSender 创建 SMTP 发送器
func NewSMTPSender(config *Config) *SMTPSender {
	return &SMTPSender{
		config: config,
	}
}

// Send 同步发送邮件
func (s *SMTPSender) Send(ctx context.Context, msg *Message) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// 构建邮件内容
	emailBody := s.buildEmail(msg)

	// 设置认证
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// 构建收件人列表
	recipients := append(msg.To, msg.Cc...)
	recipients = append(recipients, msg.Bcc...)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	var err error
	if s.config.TLS {
		// 使用 TLS（直接 TLS 连接，如 465 端口）
		tlsConfig := &tls.Config{
			ServerName: s.config.Host,
		}
		conn, connErr := tls.Dial("tcp", addr, tlsConfig)
		if connErr != nil {
			return fmt.Errorf("failed to connect: %w", connErr)
		}
		defer conn.Close()

		client, clientErr := smtp.NewClient(conn, s.config.Host)
		if clientErr != nil {
			return fmt.Errorf("failed to create client: %w", clientErr)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		from := s.config.From
		if s.config.FromName != "" {
			from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
		}

		if err = client.Mail(from); err != nil {
			return fmt.Errorf("mail command failed: %w", err)
		}

		for _, recipient := range recipients {
			if err = client.Rcpt(recipient); err != nil {
				return fmt.Errorf("rcpt command failed for %s: %w", recipient, err)
			}
		}

		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("data command failed: %w", err)
		}

		_, err = writer.Write([]byte(emailBody))
		if err != nil {
			return fmt.Errorf("write failed: %w", err)
		}

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("close failed: %w", err)
		}
	} else {
		// 使用 STARTTLS（如 587 端口）或普通连接（如 25 端口）
		// smtp.SendMail 会自动处理 STARTTLS
		err = smtp.SendMail(addr, auth, s.config.From, recipients, []byte(emailBody))
		if err != nil {
			return fmt.Errorf("send mail failed: %w", err)
		}
	}

	return nil
}

// buildEmail 构建邮件内容
func (s *SMTPSender) buildEmail(msg *Message) string {
	from := s.config.From
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	}

	headers := fmt.Sprintf("From: %s\r\n", from)
	headers += fmt.Sprintf("To: %s\r\n", joinEmails(msg.To))
	if len(msg.Cc) > 0 {
		headers += fmt.Sprintf("Cc: %s\r\n", joinEmails(msg.Cc))
	}
	headers += fmt.Sprintf("Subject: %s\r\n", msg.Subject)

	// 如果提供了 HTML，使用 multipart
	if msg.HTML != "" {
		boundary := "----=_NextPart_" + fmt.Sprintf("%d", time.Now().UnixNano())
		headers += "MIME-Version: 1.0\r\n"
		headers += fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary)
		headers += "\r\n"

		body := fmt.Sprintf("--%s\r\n", boundary)
		body += "Content-Type: text/plain; charset=UTF-8\r\n"
		body += "\r\n"
		body += msg.Body + "\r\n"
		body += fmt.Sprintf("\r\n--%s\r\n", boundary)
		body += "Content-Type: text/html; charset=UTF-8\r\n"
		body += "\r\n"
		body += msg.HTML + "\r\n"
		body += fmt.Sprintf("\r\n--%s\r\n", boundary)

		return headers + body
	}

	// 纯文本邮件
	headers += "Content-Type: text/plain; charset=UTF-8\r\n"
	headers += "\r\n"
	return headers + msg.Body
}

// joinEmails 连接邮箱地址
func joinEmails(emails []string) string {
	result := ""
	for i, email := range emails {
		if i > 0 {
			result += ", "
		}
		result += email
	}
	return result
}

// Validate 验证消息
func (m *Message) Validate() error {
	if len(m.To) == 0 {
		return fmt.Errorf("recipients (To) cannot be empty")
	}
	if m.Subject == "" {
		return fmt.Errorf("subject cannot be empty")
	}
	if m.Body == "" && m.HTML == "" {
		return fmt.Errorf("body or HTML must be provided")
	}
	return nil
}
