package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type Sender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSender(host, port, username, password, from string) *Sender {
	return &Sender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *Sender) Send(to, subject, body string) error {
	if s.host == "" || s.username == "" || s.password == "" {
		log.Printf("[EMAIL LOG MODE] To: %s | Subject: %s | Body: %s", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	message := []byte(
		"From: " + s.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, message)
}
