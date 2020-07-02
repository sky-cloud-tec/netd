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

package sg6000

import (
	"fmt"
	"io"
	"regexp"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"golang.org/x/crypto/ssh"
)

func init() {
	// register hillstone
	cli.OperatorManagerInstance.Register(`(?i)hillstone\.SG-6000-VM01\..*`, createOpHillstone())
}

type opHillstone struct {
	lineBeak    string // \r\n \n
	transitions map[string][]string
	prompts     map[string][]*regexp.Regexp
	errs        []*regexp.Regexp
}

func createOpHillstone() cli.Operator {
	loginPrompt := regexp.MustCompile(`[[:alnum:]._\-~]+(\([[:alnum:]]+\))?# ?$`)
	configurePrompt := regexp.MustCompile(`[[:alnum:]._\-~]+(\([[:alnum:]]+\))?\(config\)# ?$`)

	return &opHillstone{
		transitions: map[string][]string{
			"login->configure": {"configure"},
			"configure->login": {"exit"},
		},
		prompts: map[string][]*regexp.Regexp{
			"login":     {loginPrompt},
			"configure": {configurePrompt},
		},

		errs: []*regexp.Regexp{
			regexp.MustCompile(`\^-+incomplete command`),
			regexp.MustCompile(`\^-+unrecognized keyword\s+`),
			regexp.MustCompile(`^Error:[\s\S]+`),
		},
		lineBeak: "\n",
	}
}

func (s *opHillstone) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}

func (s *opHillstone) GetEncoding() string {
	return ""
}

func (s *opHillstone) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

func (s *opHillstone) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

func (s *opHillstone) GetLinebreak() string {
	return s.lineBeak
}

func (s *opHillstone) GetStartMode() string {
	return "login"
}

// RegisterMode ...
func (s *opHillstone) RegisterMode(req *protocol.CliRequest) error {
	return nil
}

func (s *opHillstone) GetSSHInitializer() cli.SSHInitializer {
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
		// open channel
		if err := session.Shell(); err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create shell failed, %s", err)
		}
		return r, w, session, nil
	}
}
