package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"

	"xinhang-backend/config"
)

var cfg *config.Config

func Init(c *config.Config) {
	cfg = c
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		log.Println("WARNING: SMTP not configured, email verification disabled")
	} else {
		log.Printf("SMTP configured: host=%s, user=%s, port=%s", cfg.SMTPHost, cfg.SMTPUser, cfg.SMTPPort)
	}
}

func IsEnabled() bool {
	return cfg != nil && cfg.SMTPHost != "" && cfg.SMTPUser != ""
}

func SendVerificationCode(to, code string) error {
	if !IsEnabled() {
		return fmt.Errorf("SMTP not configured")
	}

	subject := "【新航实验国际学校】注册验证码"
	body := fmt.Sprintf(`<html><body style="font-family:sans-serif;padding:20px;">
<h2 style="color:#1a73e8;">新航实验国际学校</h2>
<p>您好，您正在注册新航实验国际学校报名系统账号。</p>
<p>您的验证码为：</p>
<div style="background:#f5f5f5;padding:16px 24px;display:inline-block;border-radius:8px;margin:12px 0;">
  <span style="font-size:32px;font-weight:bold;letter-spacing:8px;color:#1a73e8;">%s</span>
</div>
<p>验证码 <strong>5 分钟</strong>内有效，请勿将验证码泄露给他人。</p>
<p style="color:#999;font-size:12px;margin-top:24px;">如非本人操作，请忽略此邮件。</p>
</body></html>`, code)

	from := cfg.SMTPUser
	displayFrom := cfg.SMTPFrom
	if displayFrom == "" {
		displayFrom = from
	}

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		displayFrom, from, to, subject, body)

	addr := net.JoinHostPort(cfg.SMTPHost, cfg.SMTPPort)

	if cfg.SMTPPort == "465" {
		return sendSSL(addr, from, to, []byte(msg))
	}
	return sendSTARTTLS(addr, from, to, []byte(msg))
}

type plainAuthDirect struct {
	user, pass string
}

func (a *plainAuthDirect) Start(server *smtp.ServerInfo) (string, []byte, error) {
	resp := []byte("\x00" + a.user + "\x00" + a.pass)
	return "PLAIN", resp, nil
}

func (a *plainAuthDirect) Next(fromServer []byte, more bool) ([]byte, error) {
	return nil, nil
}

func sendSSL(addr, from, to string, msg []byte) error {
	tlsConfig := &tls.Config{ServerName: cfg.SMTPHost}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial %s failed: %w", addr, err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("SMTP client failed: %w", err)
	}
	defer client.Quit()

	auth := &plainAuthDirect{user: cfg.SMTPUser, pass: cfg.SMTPPassword}
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return w.Close()
}

func sendSTARTTLS(addr, from, to string, msg []byte) error {
	auth := &plainAuthDirect{user: cfg.SMTPUser, pass: cfg.SMTPPassword}
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
