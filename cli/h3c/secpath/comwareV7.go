package secpath

import (
	"fmt"
	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
)

func init() {
	// register H3C SecPath comwareV7
	cli.OperatorManagerInstance.Register(`(?i)h3c\.secpath\..*`, createOpH3CV7())
}

type opH3CV7 struct {
	lineBeak     string // \r\n \n
	transitions  map[string][]string
	prompts      map[string][]*regexp.Regexp
	errs         []*regexp.Regexp
	encodingType string
}

func createOpH3CV7() cli.Operator {
	loginPrompt := regexp.MustCompile("<[-_[:alnum:][:digit:]]{0,}>$")
	systemViewPrompt := regexp.MustCompile(`\[[-_[:alnum:][:digit:]]{0,}]$`)

	return &opH3CV7{
		transitions: map[string][]string{
			"login->system_View": {"system-view"},
			"system_View->login": {"quit"},
		},
		prompts: map[string][]*regexp.Regexp{
			"login":       {loginPrompt},
			"system_View": {systemViewPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile("^ % "),
		},
		lineBeak:     "\n",
		encodingType: "GBK",
	}
}

func (s *opH3CV7) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}

func (s *opH3CV7) GetEncoding() string {
	return s.encodingType
}

func (s *opH3CV7) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

func (s *opH3CV7) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

func (s *opH3CV7) GetLinebreak() string {
	return s.lineBeak
}

func (s *opH3CV7) GetStartMode() string {
	return "login"
}

// RegisterMode ...
func (s *opH3CV7) RegisterMode(req *protocol.CliRequest) error {
	return nil
}

func (s *opH3CV7) GetSSHInitializer() cli.SSHInitializer {
	return func(c *ssh.Client, req *protocol.CliRequest) (io.Reader, io.WriteCloser, *ssh.Session, error) {
		var err error
		session, err := c.NewSession()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("new ssh session failed, %s", err)
		}
		// get stdout and stdin channel
		r, err := session.StdoutPipe()
		if err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create stdout pipe failed, %s", err)
		}
		w, err := session.StdinPipe()
		if err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create stdin pipe failed, %s", err)
		}
		modes := ssh.TerminalModes{
			ssh.ECHO: 1, // enable echoing
		}
		if err := session.RequestPty("vt100", 0, 2000, modes); err != nil {
			return nil, nil, nil, fmt.Errorf("request pty failed, %s", err)
		}
		if err := session.Shell(); err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create shell failed, %s", err)
		}
		return r, w, session, nil
	}
}
