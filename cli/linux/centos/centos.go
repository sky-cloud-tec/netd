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

package centos

import (
	"fmt"
	"io"
	"regexp"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"golang.org/x/crypto/ssh"
)

func init() {
	// register Centos centos
	cli.OperatorManagerInstance.Register(`(?i)linux\.centos\.(9|[0-9]{1,})`, createCentos())
}

//Centos struct
type Centos struct {
	lineBeak    string // \n
	transitions map[string][]string
	prompts     map[string][]*regexp.Regexp
	errs        []*regexp.Regexp
}

func createCentos() cli.Operator {
	userPrompt := regexp.MustCompile(`\[(.*){1,}@(.*){0,} .*\](#|\$) $`)
	return &Centos{
		// mode transition
		// login -> configure_terminal
		transitions: map[string][]string{
		},
		prompts: map[string][]*regexp.Regexp{
			"login":              {userPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile(".*command not found.*"),
			regexp.MustCompile(".*No such file or directory.* "),
			regexp.MustCompile(".*invalid option.*"),
		},
		lineBeak: "\n",
	}
}

//GetPrompts Centos
func (s *Centos) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}

//GetTransitions Centos
func (s *Centos) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

//GetErrPatterns Centos
func (s *Centos) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

//GetLinebreak Centos
func (s *Centos) GetLinebreak() string {
	return s.lineBeak
}

//GetStartMode Centos
func (s *Centos) GetStartMode() string {
	return "login"
}

//GetSSHInitializer Centos
func (s *Centos) GetSSHInitializer() cli.SSHInitializer {
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
