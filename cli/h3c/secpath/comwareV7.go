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
		encodingType: "GB18030",
	}
}

func (s *opH3CV7) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}

func (s *opH3CV7) SetPrompts(k string, regs []*regexp.Regexp) {
	s.prompts[k] = regs
}

func (s *opH3CV7) SetErrPatterns(regs []*regexp.Regexp) {
	s.errs = regs
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

func (s *opH3CV7) GetExcludes() []*regexp.Regexp {
	return nil
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
