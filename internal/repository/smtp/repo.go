package smtp

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
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

func (e *SMTPSender) SendRecoveryCode(ctx context.Context, email, code string) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendRecoveryCode")
	log.Infof("Sending recovery code to %s from %s via %s", email, from, config.Config.SMTP.Host)
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Домдели смена пароля")
	m.SetHeader("Reply-To", from)
	m.SetBody("text/html", fmt.Sprintf(`
        <h2>Ваш код подтверждения:</h2>
        <h1 style="color: #f08dcc">%s</h1>
        <p>Код действителен 5 минут.</p>
    `, code))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	go func() {
		if err := dialer.DialAndSend(m); err != nil {
			log.Errorf("gomail send: %s", err)
		} else {
			log.Printf("✅ Code sent to %s", email)
		}
	}()

	return nil
}

func (e *SMTPSender) SendVerificationCode(ctx context.Context, email, code string) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendVerificationCode")
	log.Infof("Sending verification code to %s from %s via %s", email, from, config.Config.SMTP.Host)
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Домдели подтверждение почты")
	m.SetHeader("Reply-To", from)
	m.SetBody("text/html", fmt.Sprintf(`
        <h2>Ваш код подтверждения:</h2>
        <h1 style="color: #f08dcc">%s</h1>
        <p>Код действителен 5 минут.</p>
    `, code))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)
	log.Printf("start")

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send: %s", err)
		}
		log.Printf("✅ Code sent to %s", email)
	}(ctx)

	return nil
}
