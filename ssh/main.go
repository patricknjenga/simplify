package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	Address    string
	Client     *ssh.Client
	Password   string
	Port       string
	PrivateKey []byte
	Timeout    time.Duration
	User       string
}

func New(address string, password string, port string, privateKey []byte, timeout time.Duration, user string) *Ssh {
	return &Ssh{address, &ssh.Client{}, password, port, privateKey, timeout, user}
}

func (s *Ssh) Dial() error {
	var (
		auth []ssh.AuthMethod
		err  error
	)
	if s.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey(s.PrivateKey)
		if err != nil {
			return err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if s.Password != "" {
		auth = append(auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			for i := range answers {
				answers[i] = s.Password
			}
			return answers, nil
		}))
	}
	s.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.Address, s.Port), &ssh.ClientConfig{
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         s.Timeout,
		User:            s.User,
	})
	return err
}

func (s *Ssh) Exec(command string) ([]byte, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return []byte{}, nil
	}
	defer session.Close()
	return session.Output(command)
}
