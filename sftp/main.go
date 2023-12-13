package sftp

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/patricknjenga/simplify/ssh"
	"github.com/pkg/sftp"
)

type Sftp struct {
	Client    *sftp.Client
	Directory string
	Ssh       ssh.Ssh
}

func New(address string, password string, port string, privateKey []byte, timeout time.Duration, user string, directory string) *Sftp {
	return &Sftp{&sftp.Client{}, directory, *ssh.New(address, password, port, privateKey, timeout, user)}
}

func (s *Sftp) Dial() error {
	err := s.Ssh.Dial()
	if err != nil {
		return err
	}
	s.Client, err = sftp.NewClient(s.Ssh.Client)
	return err
}

func (s *Sftp) Get(p string) ([]byte, error) {
	reader, err := s.Client.Open(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return []byte{}, err
	}
	return io.ReadAll(reader)
}

func (s *Sftp) Ls(r string) ([]os.FileInfo, error) {
	var result []os.FileInfo
	regexp, err := regexp.Compile(r)
	if err != nil {
		return result, err
	}
	fileInfos, err := s.Client.ReadDir(s.Directory)
	if err != nil {
		return result, err
	}
	for _, v := range fileInfos {
		if regexp.Match([]byte(v.Name())) {
			result = append(result, v)
		}
	}
	return result, err
}

func (s *Sftp) Mkdir(d string) error {
	return s.Client.MkdirAll(fmt.Sprintf("%s/%s", s.Directory, d))
}

func (s *Sftp) Put(p string, data []byte) (int, error) {
	writer, err := s.Client.Create(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return 0, err
	}
	return writer.Write(data)
}
