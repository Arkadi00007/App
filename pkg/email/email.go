package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSender(host string, port int, username, password, from string) *Sender {
	return &Sender{host, port, username, password, from}
}

func (s *Sender) SendVerificationCode(toEmail, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Подтверждение email — Harcer")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Добро пожаловать в Harcer!</h2>
		<p>Ваш код подтверждения:</p>
		<h1 style="letter-spacing: 4px">%s</h1>
		<p>Код действителен 15 минут.</p>
	`, code))

	return s.send(m)
}

func (s *Sender) SendResetPasswordCode(toEmail, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Сброс пароля — Harcer")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Сброс пароля</h2>
		<p>Ваш код для сброса пароля:</p>
		<h1 style="letter-spacing: 4px">%s</h1>
		<p>Код действителен 15 минут.</p>
		<p>Если вы не запрашивали сброс пароля — проигнорируйте это письмо.</p>
	`, code))

	return s.send(m)
}

func (s *Sender) send(m *gomail.Message) error {
	d := gomail.NewDialer(s.host, s.port, s.username, s.password)
	return d.DialAndSend(m)
}
