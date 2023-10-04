package smb

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"

	"github.com/hirochachacha/go-smb2"
)

type Smb struct {
	Address   string
	Directory string
	Mount     string
	Password  string
	Port      string
	Session   *smb2.Session
	Share     *smb2.Share
	User      string
}

func New(address string, directory string, mount string, password string, port string, user string) *Smb {
	return &Smb{address, directory, mount, password, port, &smb2.Session{}, &smb2.Share{}, user}
}

func (s *Smb) Dial() error {
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.Address, s.Port))
	if err != nil {
		return err
	}
	s.Session, err = (&smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.User,
			Password: s.Password,
		},
	}).Dial(connection)
	if err != nil {
		return err
	}
	s.Share, err = s.Session.Mount(s.Mount)
	return err
}

func (s *Smb) Get(p string) ([]byte, error) {
	reader, err := s.Share.Open(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return []byte{}, err
	}
	return io.ReadAll(reader)
}

func (s *Smb) Ls(r string) ([]os.FileInfo, error) {
	var result []os.FileInfo
	regexp, err := regexp.Compile(r)
	if err != nil {
		return result, err
	}
	fileInfos, err := s.Share.ReadDir(s.Directory)
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

func (s *Smb) Mkdir(d string) error {
	return s.Share.MkdirAll(fmt.Sprintf("%s/%s", s.Directory, d), 0777)
}

func (s *Smb) Put(path string, data []byte) (int, error) {
	writer, err := s.Share.Create(fmt.Sprintf("%s/%s", s.Directory, path))
	if err != nil {
		return 0, err
	}
	return writer.Write(data)
}
