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

func (e *SMTPSender) SendPromotionExpier(ctx context.Context, email string) error {
	from := config.Config.SMTP.Email
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendAnswer")
	log.Infof("Sending to %s from %s via %s", email, from, config.Config.SMTP.Host)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Не забудте продлить буст")
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", `
		<h2>Ваше платное продвижение объявления скоро истечет</h2>
		<a href="https://dom-deli.ru/my-posters" style="color: #f08dcc">Продлить буст...</a> 
	`)

	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)

	go func(ctx context.Context) {
		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("gomail send answer: %s", err)
			return
		}
		log.Printf("✅ Email sent to %s", email)
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
	<!doctype html>
	<html lang="ru">
	<body style="margin:0; padding:0; background:#f4f6fb;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6fb; padding:16px;">
			<tr>
				<td align="center">
					<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; background:#ffffff; border-radius:18px; overflow:hidden; font-family:Arial, sans-serif; box-shadow:0 8px 28px rgba(0,0,0,0.08);">

						<tr>
							<td style="background:rgb(240, 141, 204); padding:28px 32px;">
								<div style="font-size:24px; font-weight:700; color:#ffffff;">DomDeli</div>
								<div style="margin-top:6px; font-size:15px; color:#eaf3ff;">Новая заявка в соседи</div>
							</td>
						</tr>

						<tr>
							<td style="padding:32px;">
								<h2 style="margin:0 0 16px; font-size:26px; color:#1f2933;">
									У вас новая заявка 👋
								</h2>

								<p style="margin:0 0 20px; font-size:16px; line-height:26px; color:#52616b;">
									<b style="color:#1f2933;">%s %s</b> хочет добавить вас в соседи.
								</p>

								<div style="background:rgb(246, 246, 245); border:1px solid #e8eef6; border-radius:14px; padding:18px 20px; margin:22px 0;">
									<p style="margin:0; font-size:15px; line-height:24px; color:#52616b;">
										Посмотрите объявление и перейдите в заявки, чтобы ответить пользователю.
									</p>
								</div>

								<table cellpadding="0" cellspacing="0" style="margin:24px 0 14px;">
									<tr>
										<td style="background:rgb(246, 246, 245); border-radius:12px;">
											<a href="%s" target="_blank"
											style="display:inline-block; padding:14px 22px; color:rgb(26, 26, 26); text-decoration:none; font-size:15px; font-weight:700;">
												Посмотреть объявление
											</a>
										</td>
									</tr>
								</table>

								<table cellpadding="0" cellspacing="0" style="margin:0 0 22px;">
									<tr>
										<td style="background:rgb(240, 141, 204); border-radius:12px;">
											<a href="https://dom-deli.ru/friends?tab=requests" target="_blank"
											style="display:inline-block; padding:14px 22px; color:#ffffff; text-decoration:none; font-size:15px; font-weight:700;">
												Открыть заявки в соседи
											</a>
										</td>
									</tr>
								</table>

								<p style="margin:0; font-size:13px; line-height:20px; color:#8a94a6;">
									Если кнопка не открывается, скопируйте ссылку в браузер:<br>
									<a href="%s" style="color:#2f80ed;">%s</a>
								</p>
							</td>
						</tr>

						<tr>
							<td style="padding:20px 32px; background:#f8fafd; border-top:1px solid #e8eef6;">
								<p style="margin:0; font-size:13px; line-height:20px; color:#7b8794;">
									Это автоматическое письмо от DomDeli. Если вы не ожидали это уведомление, просто проигнорируйте его.
								</p>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, firstName, lastName, posterURL, posterURL, posterURL))

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
	log.Infof("SendRoommateContactsForRequester: %s, %s, %s, %s, %s, %s", email, roommateFirstName, roommateLastName, roommateEmail, roommatePhone, posterAlias)

	posterURL := fmt.Sprintf("https://dom-deli.ru/posters/%s", posterAlias)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Контакты сожителя")
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", fmt.Sprintf(`
	<!doctype html>
	<html lang="ru">
	<body style="margin:0; padding:0; background:#f4f6fb;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6fb; padding:16px;">
			<tr>
				<td align="center">
					<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; background:#ffffff; border-radius:18px; overflow:hidden; font-family:Arial, sans-serif; box-shadow:0 8px 28px rgba(0,0,0,0.08);">

						<tr>
							<td style="background:rgb(240, 141, 204); padding:28px 32px;">
								<div style="font-size:24px; font-weight:700; color:#ffffff;">DomDeli</div>
								<div style="margin-top:6px; font-size:15px; color:#eaf3ff;">Контакты сожителя</div>
							</td>
						</tr>

						<tr>
							<td style="padding:32px;">
								<h2 style="margin:0 0 16px; font-size:26px; color:#1f2933;">
									Контакты сожителя 🎉
								</h2>

								<p style="margin:0 0 20px; font-size:16px; line-height:26px; color:#52616b;">
									Вы добавили <b style="color:#1f2933;">%s %s</b> в соседи.
									Теперь вы можете связаться с пользователем напрямую.
								</p>

								<div style="background:rgb(246, 246, 245); border:1px solid #e8eef6; border-radius:14px; padding:18px 20px; margin:22px 0;">
									<table width="100%%" cellpadding="0" cellspacing="0">
										<tr>
											<td style="padding:8px 0; font-size:14px; color:#7b8794;">Email</td>
											<td style="padding:8px 0; font-size:15px; color:#1f2933; font-weight:700;">
												%s
											</td>
										</tr>
										<tr>
											<td style="padding:8px 0; font-size:14px; color:#7b8794;">Тел.</td>
											<td style="padding:8px 0; font-size:15px; color:#1f2933; font-weight:700;">
												%s
											</td>
										</tr>
									</table>
								</div>

								<table cellpadding="0" cellspacing="0" style="margin:24px 0 14px;">
									<tr>
										<td style="background:rgb(246, 246, 245); border-radius:12px;">
											<a href="%s" target="_blank"
											style="display:inline-block; padding:14px 22px; color:rgb(26, 26, 26); text-decoration:none; font-size:15px; font-weight:700;">
												Посмотреть объявление
											</a>
										</td>
									</tr>
								</table>

								<table cellpadding="0" cellspacing="0" style="margin:0 0 22px;">
									<tr>
										<td style="background:rgb(240, 141, 204); border-radius:12px;">
											<a href="https://dom-deli.ru/friends" target="_blank"
											style="display:inline-block; padding:14px 22px; color:#ffffff; text-decoration:none; font-size:15px; font-weight:700;">
												Открыть список соседей
											</a>
										</td>
									</tr>
								</table>

								<p style="margin:0; font-size:13px; line-height:20px; color:#8a94a6;">
									Если кнопка не открывается, скопируйте ссылку в браузер:<br>
									<a href="https://dom-deli.ru/friends" style="color:#2f80ed;">https://dom-deli.ru/friends</a>
								</p>
							</td>
						</tr>

						<tr>
							<td style="padding:20px 32px; background:#f8fafd; border-top:1px solid #e8eef6;">
								<p style="margin:0; font-size:13px; line-height:20px; color:#7b8794;">
									Это автоматическое письмо от DomDeli. Если вы не ожидали это уведомление, просто проигнорируйте его.
								</p>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, roommateFirstName, roommateLastName, roommateEmail, roommatePhone, posterURL))

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
	log.Infof("SendRoommateContactsForAccepted: %s, %s, %s, %s, %s, %s", email, roommateFirstName, roommateLastName, roommateEmail, roommatePhone, posterAlias)

	posterURL := fmt.Sprintf("https://dom-deli.ru/posters/%s", posterAlias)

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"DomDeli" <%s>`, from))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Вашу заявку в соседи приняли")
	m.SetHeader("Reply-To", from)

	m.SetBody("text/html", fmt.Sprintf(`
		<!doctype html>
			<html lang="ru">
			<body style="margin:0; padding:0; background:#f4f6fb;">
				<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6fb; padding:16px;">
					<tr>
						<td align="center">
							<table width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; background:#ffffff; border-radius:18px; overflow:hidden; font-family:Arial, sans-serif; box-shadow:0 8px 28px rgba(0,0,0,0.08);">

								<tr>
									<td style="background:rgb(240, 141, 204); padding:28px 32px;">
										<div style="font-size:24px; font-weight:700; color:#ffffff;">DomDeli</div>
										<div style="margin-top:6px; font-size:15px; color:#eaf3ff;">Заявку приняли</div>
									</td>
								</tr>

								<tr>
									<td style="padding:32px;">
										<h2 style="margin:0 0 16px; font-size:26px; color:#1f2933;">
											Вашу заявку приняли 🎉
										</h2>

										<p style="margin:0 0 20px; font-size:16px; line-height:26px; color:#52616b;">
											<b style="color:#1f2933;">%s %s</b> добавил вас в соседи.
											Теперь вы можете связаться с пользователем напрямую.
										</p>

										<div style="background:rgb(246, 246, 245); border:1px solid #e8eef6; border-radius:14px; padding:18px 20px; margin:22px 0;">
											<table width="100%%" cellpadding="0" cellspacing="0">
												<tr>
													<td style="padding:8px 0; font-size:14px; color:#7b8794;">Email</td>
													<td style="padding:8px 0; font-size:15px; color:#1f2933; font-weight:700;">
														%s
													</td>
												</tr>
												<tr>
													<td style="padding:8px 0; font-size:14px; color:#7b8794;">Тел.</td>
													<td style="padding:8px 0; font-size:15px; color:#1f2933; font-weight:700;">
														%s
													</td>
												</tr>
											</table>
										</div>

										<table cellpadding="0" cellspacing="0" style="margin:24px 0 14px;">
											<tr>
												<td style="background:rgb(246, 246, 245); border-radius:12px;">
													<a href="%s" target="_blank"
													style="display:inline-block; padding:14px 22px; color:rgb(26, 26, 26); text-decoration:none; font-size:15px; font-weight:700;">
														Посмотреть объявление
													</a>
												</td>
											</tr>
										</table>

										<table cellpadding="0" cellspacing="0" style="margin:0 0 22px;">
											<tr>
												<td style="background:rgb(240, 141, 204); border-radius:12px;">
													<a href="https://dom-deli.ru/friends" target="_blank"
													style="display:inline-block; padding:14px 22px; color:#ffffff; text-decoration:none; font-size:15px; font-weight:700;">
														Открыть список соседей
													</a>
												</td>
											</tr>
										</table>

										<p style="margin:0; font-size:13px; line-height:20px; color:#8a94a6;">
											Если кнопка не открывается, скопируйте ссылку в браузер:<br>
											<a href="https://dom-deli.ru/friends" style="color:#2f80ed;">https://dom-deli.ru/friends</a>
										</p>
									</td>
								</tr>

								<tr>
									<td style="padding:20px 32px; background:#f8fafd; border-top:1px solid #e8eef6;">
										<p style="margin:0; font-size:13px; line-height:20px; color:#7b8794;">
											Это автоматическое письмо от DomDeli. Если вы не ожидали это уведомление, просто проигнорируйте его.
										</p>
									</td>
								</tr>

							</table>
						</td>
					</tr>
				</table>
			</body>
			</html>
		`, roommateFirstName, roommateLastName, roommateEmail, roommatePhone, posterURL))

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
