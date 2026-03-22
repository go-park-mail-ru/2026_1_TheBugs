package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
)

type SMTPSender struct {
	cred    *smtp.Auth
	address string
}

func NewSMTPSender(addres string, cred *smtp.Auth) *SMTPSender {
	return &SMTPSender{
		cred:    cred,
		address: addres,
	}
}
func (e SMTPSender) SendCode(ctx context.Context, email string, code string) error {
	from := config.Config.SMTP.Email
	to := []string{email}

	msg := []byte(
		fmt.Sprintf("From: %s\r\n", from) +
			fmt.Sprintf("To: %s\r\n", email) +
			"Subject: DomDeli verification code\r\n" +
			"\r\n" +
			fmt.Sprintf("Your code: %s", code))

	err := smtp.SendMail(e.address, nil, from, to, msg)

	if err != nil {
		return fmt.Errorf("smtp.SendMail: %e", err)
	}

	return nil
}
