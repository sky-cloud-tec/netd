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

package ios

import (
	"fmt"
	"io"
	"regexp"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"golang.org/x/crypto/ssh"
)

func init() {
	// register switch ios
	cli.OperatorManagerInstance.Register(`(?i)cisco\.ios\..*`, createSwitchIos())
}

//SwitchIos struct
type SwitchIos struct {
	lineBeak    string // \r\n \n
	transitions map[string][]string
	prompts     map[string][]*regexp.Regexp
	errs        []*regexp.Regexp
}

func createSwitchIos() cli.Operator {
	loginPrompt := regexp.MustCompile("^[[:alnum:]._-]+> ?$")
	loginEnablePrompt := regexp.MustCompile("[[:alnum:]]{1,}(-[[:alnum:]]+){0,}#$")
	configTerminalPrompt := regexp.MustCompile(`[[:alnum:]]{1,}(-[[:alnum:]]+){0,}\(config\)#$`)
	return &SwitchIos{
		// mode transition
		// login_enable -> configure_terminal
		transitions: map[string][]string{
			"login_enable->configure_terminal": {"config terminal"},
			"configure_terminal->login_enable": {"exit"},
		},
		prompts: map[string][]*regexp.Regexp{
			"login_or_login_enable": {loginPrompt, loginEnablePrompt},
			"login":                 {loginPrompt},
			"login_enable":          {loginEnablePrompt},
			"configure_terminal":    {configTerminalPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile("^Command authorization failed\\.$"),
			regexp.MustCompile("^% "),
			regexp.MustCompile("^Command rejected:"),
		},
		lineBeak: "\n",
	}
}

//GetPrompts SwitchIos
func (s *SwitchIos) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}

// GetEncoding return device encoding type as string
func (s *SwitchIos) GetEncoding() string {
	return ""
}

// GetExcludes return excluded prommpt pattern
func (s *SwitchIos) GetExcludes() []*regexp.Regexp {
	return nil
}

//GetTransitions SwitchIos
func (s *SwitchIos) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

//GetErrPatterns SwitchIos
func (s *SwitchIos) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

//GetLinebreak SwitchIos
func (s *SwitchIos) GetLinebreak() string {
	return s.lineBeak
}

//GetStartMode SwitchIos
func (s *SwitchIos) GetStartMode() string {
	return "login_or_login_enable"
}

// RegisterMode ...
func (s *SwitchIos) RegisterMode(req *protocol.CliRequest) error {
	return nil
}

//GetSSHInitializer SwitchIos
func (s *SwitchIos) GetSSHInitializer() cli.SSHInitializer {
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
