package brocadeswitch

import (
	"fmt"
	"io"
	"regexp"

	"github.com/sky-cloud-tec/netd/cli"
	"golang.org/x/crypto/ssh"
)

func init() {
	// register asa 9.x+
	cli.OperatorManagerInstance.Register(`(?i)brocade\.brocadeswitch\..*`, createbrocadeSwitch())
}

type brocadeSwitch struct {
	lineBeak    string // \r\n \n
	transitions map[string][]string
	prompts     map[string][]*regexp.Regexp
	errs        []*regexp.Regexp
}

func createbrocadeSwitch() cli.Operator {
	loginPrompt := regexp.MustCompile("Lab_610_1:admin> ")
	return &brocadeSwitch{
		// mode transition
		// login_enable -> configure_terminal
		transitions: map[string][]string{
		},
		prompts: map[string][]*regexp.Regexp{
			"login":                 {loginPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile("^ERROR: "),
		},
		lineBeak: "\n",
	}
}

func (s *brocadeSwitch) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}
func (s *brocadeSwitch) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

func (s *brocadeSwitch) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

func (s *brocadeSwitch) GetLinebreak() string {
	return s.lineBeak
}

func (s *brocadeSwitch) GetStartMode() string {
	return "login"
}

func (s *brocadeSwitch) GetSSHInitializer() cli.SSHInitializer {
	return func(c *ssh.Client) (io.Reader, io.WriteCloser, *ssh.Session, error) {
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
		if err := session.Shell(); err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create shell failed, %s", err)
		}
		return r, w, session, nil
	}
}
