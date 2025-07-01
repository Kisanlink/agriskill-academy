// File: internal/notification/service.go

package notification

import (
	"os"

	"gopkg.in/gomail.v2"
)

type NotificationService interface {
	SendEmail(to, subject, body string) error
}

type mailService struct {
	From     string
	SMTPHost string
	SMTPPort int
	Password string
}

func NewMailService() NotificationService {
	return &mailService{
		From:     os.Getenv("MAIL_FROM"),
		SMTPHost: os.Getenv("MAIL_HOST"),
		SMTPPort: 587, // or os.Getenv("MAIL_PORT")
		Password: os.Getenv("MAIL_PASS"),
	}
}

func (s *mailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.From, s.Password)
	return d.DialAndSend(m)
}
