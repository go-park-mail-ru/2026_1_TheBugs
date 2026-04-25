package smtp

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

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send: %s", err)
		}
		log.Printf("✅ Code sent to %s", email)
	}(ctx)

	return nil
}

func (e *SMTPSender) SendAnswer(ctx context.Context, email string, orderID int, answer string) error {
	from := config.Config.SMTP.Email

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("Ответ поддержки по обращению #%d", orderID))
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Ответ поддержки</h2>
		<p><b>Номер обращения:</b> #%d</p>
		<p>%s</p>
	`, orderID, answer))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send answer: %s", err)
			return
		}
		log.Printf("✅ Answer sent to %s", email)
	}(ctx)

	return nil
}
