// NetD makes network device operations easy.
// Copyright (C) 2019  sky-cloud.net
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package asa

import (
	"fmt"
	"io"
	"regexp"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"golang.org/x/crypto/ssh"
)

func init() {
	// register asa 9.x+
	cli.OperatorManagerInstance.Register(`(?i)cisco\.asa[a-z]{0,}\.(9|[0-9]{1,})\..*`, createOp9xPlus())
}

type op9xPlus struct {
	lineBeak    string // \r\n \n
	transitions map[string][]string
	prompts     map[string][]*regexp.Regexp
	errs        []*regexp.Regexp
}

func createOp9xPlus() cli.Operator {
	loginPrompt := regexp.MustCompile("[[:alnum:]]{1,}(-[[:alnum:]]+){0,}> $")
	loginEnablePrompt := regexp.MustCompile("[[:alnum:]]{1,}(-[[:alnum:]]+){0,}# $")
	configTerminalPrompt := regexp.MustCompile(`[[:alnum:]]{1,}(-[[:alnum:]]+){0,}\(config\)# $`)
	return &op9xPlus{
		// mode transition
		// login_enable -> configure_terminal
		transitions: map[string][]string{
			"login_enable->configure_terminal": {"configure terminal"},
			"configure_terminal->login_enable": {"exit"},
		},
		prompts: map[string][]*regexp.Regexp{
			"login_or_login_enable": {loginPrompt, loginEnablePrompt},
			"login":                 {loginPrompt},
			"login_enable":          {loginEnablePrompt},
			"configure_terminal":    {configTerminalPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile("^ERROR: .*$"),
		},
		lineBeak: "\n",
	}
}

func (s *op9xPlus) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}
func (s *op9xPlus) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

func (s *op9xPlus) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

func (s *op9xPlus) GetLinebreak() string {
	return s.lineBeak
}

func (s *op9xPlus) GetStartMode() string {
	return "login_or_login_enable"
}

func (s *op9xPlus) RegisterMode(req *protocol.CliRequest) error {
	return nil
}

func (s *op9xPlus) GetSSHInitializer() cli.SSHInitializer {
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
		if err := session.Shell(); err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create shell failed, %s", err)
		}
		return r, w, session, nil
	}
}
