// File: pkg/mailutil/mailutil.go

package mailutil

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	From     string
	Password string
	Host     string
	Port     int
}

func NewMailer() *Mailer {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return &Mailer{
		From:     os.Getenv("MAIL_FROM"),
		Password: os.Getenv("MAIL_PASS"),
		Host:     os.Getenv("MAIL_HOST"),
		Port:     port,
	}
}

func (m *Mailer) SendMail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer(m.Host, m.Port, m.From, m.Password)
	return d.DialAndSend(msg)
}
