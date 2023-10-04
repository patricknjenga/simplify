package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	Address    string
	Client     *ssh.Client
	Password   string
	Port       string
	PrivateKey []byte
	User       string
}

func New(address string, password string, port string, privateKey []byte, user string) *Ssh {
	return &Ssh{address, &ssh.Client{}, password, port, privateKey, user}
}

func (s *Ssh) Dial() error {
	signer, err := ssh.ParsePrivateKey(s.PrivateKey)
	if err != nil && s.PrivateKey != nil {
		return err
	}
	s.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.Address, s.Port), &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = s.Password
				}
				return answers, nil
			}),
			ssh.Password(s.Password),
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            s.User,
	})
	return err
}

func (s *Ssh) Exec(command string) ([]byte, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return []byte{}, nil
	}
	return session.Output(command)
}
