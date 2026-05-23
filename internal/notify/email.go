package notify

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailAlerter sends alert notifications via SMTP.
type EmailAlerter struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       []string
}

// NewEmailAlerter constructs an EmailAlerter from the given config.
func NewEmailAlerter(host, port, username, password, from string, to []string) *EmailAlerter {
	return &EmailAlerter{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}
}

// Alert sends an email notification for an overdue job.
func (e *EmailAlerter) Alert(jobName string) error {
	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	subject := fmt.Sprintf("Subject: [cronitor-local] Job overdue: %s\r\n", jobName)
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\n", e.from, strings.Join(e.to, ", "))
	body := fmt.Sprintf("\r\nAlert: cron job '%s' has not run within its expected schedule.\r\n", jobName)
	msg := []byte(subject + headers + body)

	return smtp.SendMail(addr, auth, e.from, e.to, msg)
}
