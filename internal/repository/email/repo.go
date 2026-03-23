package email

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	gomail "gopkg.in/gomail.v2"
)

type SMTPSender struct {
	username string
	password string
	host     string
	port     int
}

func NewSMTPSender(host string, port int, username, password string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (e *SMTPSender) SendCode(ctx context.Context, email, code string) error {
	from := config.Config.SMTP.Email

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Домдели смена пароля")
	m.SetHeader("Reply-To", from)
	m.SetBody("text/html", fmt.Sprintf(`
        <h2>Ваш код подтверждения:</h2>
        <h1 style="color: #007bff">%s</h1>
        <p>Код действителен 5 минут.</p>
    `, code))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("gomail send: %w", err)
	}

	log.Printf("✅ Code sent to %s", email)
	return nil
}
