package notify

import (
	"fmt"
	"net/smtp"
)

// EmailConfig holds SMTP configuration for sending alert emails.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
}

// EmailAlerter sends alert notifications via email.
type EmailAlerter struct {
	cfg EmailConfig
}

// NewEmailAlerter creates a new EmailAlerter with the given config.
func NewEmailAlerter(cfg EmailConfig) *EmailAlerter {
	return &EmailAlerter{cfg: cfg}
}

// Alert sends an email notification for the given job name and reason.
func (e *EmailAlerter) Alert(jobName, reason string) error {
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)

	subject := fmt.Sprintf("[cronitor-local] Alert: %s", jobName)
	body := fmt.Sprintf("Job: %s\nReason: %s", jobName, reason)
	msg := []byte(fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		e.cfg.To, e.cfg.From, subject, body))

	return smtp.SendMail(addr, auth, e.cfg.From, []string{e.cfg.To}, msg)
}
