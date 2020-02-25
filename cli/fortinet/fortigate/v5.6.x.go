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

package fortigate

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"github.com/songtianyi/rrframework/logs"
	"golang.org/x/crypto/ssh"
)

type opFortinet struct {
	lineBreak   string // /r/n \n
	transitions map[string][]string
	prompts     map[string][]*regexp.Regexp
	errs        []*regexp.Regexp
}

func init() {
	cli.OperatorManagerInstance.Register(`(?i)fortinet\.FortiGate-VM64-KVM\..*`, createOpfortinet())
}

func createOpfortinet() cli.Operator {
	loginPrompt := regexp.MustCompile(`[[:alnum:]]{1,}[[:alnum:]-_]{0,} # $`)
	return &opFortinet{
		transitions: map[string][]string{},
		prompts: map[string][]*regexp.Regexp{
			"login": {loginPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile("^Unknown action 0$"),
			regexp.MustCompile("^command parse error"),
			regexp.MustCompile("^value parse error"),
			regexp.MustCompile("^Command fail. Return code"),
			regexp.MustCompile("^please use 'end' to return to root shell"),
			regexp.MustCompile("^entry not found in datasource"),
			regexp.MustCompile("^node_check_object fail"),
		},
		lineBreak: "\n",
	}
}

func (s *opFortinet) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.prompts[k]; ok {
		return v
	}
	return nil
}

func (s *opFortinet) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

func (s *opFortinet) GetErrPatterns() []*regexp.Regexp {
	return s.errs
}

func (s *opFortinet) GetStartMode() string {
	return "login"
}

func (s *opFortinet) GetLinebreak() string {
	return s.lineBreak
}

func (s *opFortinet) RegisterMode(req *protocol.CliRequest) error {
	if s.GetPrompts(req.Mode) != nil {
		return nil
	}
	// no pattern for this mode
	// try insert
	logs.Info(req.LogPrefix, "registering pattern for mode", req.Mode)
	if strings.EqualFold(req.Mode, "global") {
		// global vdom
		s.prompts[req.Mode] = []*regexp.Regexp{
			regexp.MustCompile(`[[:alnum:]]{1,}[[:alnum:]-_]{0,} \(` + req.Mode + `\) # $`),
		}
		s.transitions["login->"+req.Mode] = []string{"config global"}
		s.transitions[req.Mode+"->"+"login"] = []string{"end"}
	} else {
		s.prompts[req.Mode] = []*regexp.Regexp{
			regexp.MustCompile(`[[:alnum:]]{1,}[[:alnum:]-_]{0,} \(` + req.Mode + `\) # $`),
		}
		s.transitions["login->"+req.Mode] = []string{"config vdom\n\t" +
			"edit " + req.Mode +
			``}
		s.transitions[req.Mode+"->"+"login"] = []string{"end"}
		// no matter global registered or not, register transition
		s.transitions["global->"+req.Mode] = []string{"end\nconfig vdom\n\t" +
			"edit " + req.Mode +
			``}
		s.transitions[req.Mode+"->global"] = []string{"end\nconfig global"}
	}
	logs.Debug(req.LogPrefix, s)
	return nil
}

func (s *opFortinet) GetSSHInitializer() cli.SSHInitializer {
	return func(c *ssh.Client, req *protocol.CliRequest) (io.Reader, io.WriteCloser, *ssh.Session, error) {
		if err := s.RegisterMode(req); err != nil {
			return nil, nil, nil, err
		}
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
