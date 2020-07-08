package mail

import (
	"fmt"
	"net"
	"net/smtp"
)

type Mail struct {
	server   string
	username string
	auth     smtp.Auth
}

func New(host string, port string, username string, password string) *Mail {
	return &Mail{
		server:   net.JoinHostPort(host, port),
		username: username,
		auth:     smtp.PlainAuth("", username, password, host),
	}
}

func (m *Mail) SendVerificationEmail(to string, link string) error {
	// TODO more domain-specific stuff
	return m.sendEmail(to, "Verify your exlibris account", fmt.Sprintf("Thank you for registering on exlibris! To verify your account, visit this link:\r\n\r\n%s", link))
}

func (m *Mail) sendEmail(to string, subject string, body string) error {
	recipients := []string{to}
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n\r\n%s",
		to, subject, body,
	))
	return smtp.SendMail(m.server, m.auth, m.username, recipients, msg)
}
