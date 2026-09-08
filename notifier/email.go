package notifier

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends an email using the configured SMTP server.
func SendEmail(subject string, body string, toAddress string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || password == "" || from == "" {
		return fmt.Errorf("SMTP environment variables are not fully configured")
	}

	auth := smtp.PlainAuth("", user, password, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", toAddress, from, subject, body))

	err := smtp.SendMail(addr, auth, from, []string{toAddress}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
