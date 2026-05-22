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

func (e *SMTPSender) SendAnswer(ctx context.Context, email string, orderID int, answer string) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendAnswer")
	log.Infof("Sending verification code to %s from %s via %s", email, from, config.Config.SMTP.Host)

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

func (e *SMTPSender) SendRoommateMatch(ctx context.Context, email string, firstName string, lastName string, posterAlias string) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendRoommateMatch")
	log.Infof("Sending roommate match notification to %s from %s via %s", email, from, config.Config.SMTP.Host)

	posterURL := fmt.Sprintf("https://dom-deli.ru/posters/%s", posterAlias)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Вас хотят добавить в соседи")
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", fmt.Sprintf(`
		<h2>У вас новая заявка</h2>

		<p><b>%s %s</b> хочет добавить вас в соседи.</p>

		<p><b>Объявление:</b> <a href="%s">%s</a></p>

		<p>
			Зайдите в DomDeli, чтобы посмотреть анкету и ответить взаимностью:
			<a href="https://dom-deli.ru/profile">https://dom-deli.ru/profile</a>
		</p>
	`, firstName, lastName, posterURL, posterURL))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send roommate match: %s", err)
			return
		}
		log.Printf("✅ Roommate match notification sent to %s", email)
	}(ctx)

	return nil
}

func (e *SMTPSender) SendRoommateContactsForRequester(
	ctx context.Context,
	email string,
	roommateFirstName string,
	roommateLastName string,
	roommateEmail string,
	roommatePhone string,
	posterAlias string,
) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendRoommateContactsForRequester")
	log.Infof("Sending roommate contacts to %s from %s via %s", email, from, config.Config.SMTP.Host)

	posterURL := fmt.Sprintf("https://dom-deli.ru/posters/%s", posterAlias)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Контакты сожителя")
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Контакты сожителя</h2>

		<p>Вы добавили <b>%s %s</b> в соседи.</p>

		<p><b>Объявление:</b> <a href="%s">%s</a></p>

		<p>Теперь вы можете связаться с пользователем:</p>

		<p><b>Email:</b> %s</p>
		<p><b>Телефон:</b> %s</p>

		<p>
			Список ваших соседей вы можете посмотреть в профиле:
			<a href="https://dom-deli.ru/profile">https://dom-deli.ru/profile</a>
		</p>
	`, roommateFirstName, roommateLastName, posterURL, posterURL, roommateEmail, roommatePhone))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send roommate contacts for requester: %s", err)
			return
		}
		log.Printf("✅ Roommate contacts sent to %s", email)
	}(ctx)

	return nil
}

func (e *SMTPSender) SendRoommateContactsForAccepted(
	ctx context.Context,
	email string,
	roommateFirstName string,
	roommateLastName string,
	roommateEmail string,
	roommatePhone string,
	posterAlias string,
) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendRoommateContactsForAccepted")
	log.Infof("Sending accepted roommate contacts to %s from %s via %s", email, from, config.Config.SMTP.Host)

	posterURL := fmt.Sprintf("https://dom-deli.ru/posters/%s", posterAlias)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Вашу заявку в соседи приняли")
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Вашу заявку в соседи приняли</h2>

		<p><b>%s %s</b> добавил вас в соседи.</p>

		<p><b>Объявление:</b> <a href="%s">%s</a></p>

		<p>Теперь вы можете связаться с пользователем:</p>

		<p><b>Email:</b> %s</p>
		<p><b>Телефон:</b> %s</p>

		<p>
			Список ваших соседей вы можете посмотреть в профиле:
			<a href="https://dom-deli.ru/profile">https://dom-deli.ru/profile</a>
		</p>
	`, roommateFirstName, roommateLastName, posterURL, posterURL, roommateEmail, roommatePhone))

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send accepted roommate contacts: %s", err)
			return
		}
		log.Printf("✅ Accepted roommate contacts sent to %s", email)
	}(ctx)

	return nil
}
