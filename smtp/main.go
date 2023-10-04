package smtp

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Host string
	Pass string
	Port int
	User string
}

func New(host string, pass string, port int, user string) *Smtp {
	return &Smtp{host, pass, port, user}
}

func (s *Smtp) Send(attach []string, body string, contentType string, embed []string, from []string, subject []string, to []string) error {
	mail := gomail.NewMessage()
	mail.SetBody(contentType, body)
	mail.SetHeaders(map[string][]string{
		"From":    from,
		"Subject": subject,
		"To":      to,
	})
	for _, v := range attach {
		mail.Attach(v)
	}
	for _, v := range embed {
		mail.Embed(v)
	}
	dialer := gomail.NewDialer(s.Host, s.Port, s.User, s.Pass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialer.DialAndSend(mail)
}
