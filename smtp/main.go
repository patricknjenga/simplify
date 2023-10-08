package smtp

import (
	"crypto/tls"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Host string
	Pass string
	Port string
	User string
}

func New(host string, pass string, port string, user string) *Smtp {
	return &Smtp{host, pass, port, user}
}

func (s *Smtp) Send(attach []string, body string, contentType string, embed []string, from []string, subject []string, to []string) error {
	port, err := strconv.Atoi(s.Port)
	if err != nil {
		return err
	}
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
	dialer := gomail.NewDialer(s.Host, port, s.User, s.Pass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialer.DialAndSend(mail)
}
