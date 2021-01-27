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
	loginPrompt := regexp.MustCompile(`[[:alnum:]]{1,}[[:alnum:]-_]{0,} (#|\$) $`)
	return &opFortinet{
		transitions: map[string][]string{},
		prompts: map[string][]*regexp.Regexp{
			"login": {loginPrompt},
		},
		errs: []*regexp.Regexp{
			regexp.MustCompile("^Unknown action 0$"),
			regexp.MustCompile(" # Unknown action 0$"),
			regexp.MustCompile("^command parse error"),
			regexp.MustCompile("^value parse error"),
			regexp.MustCompile("^Command fail. Return code"),
			regexp.MustCompile(" # Command fail. Return code"),
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

func (s *opFortinet) SetPrompts(k string, regs []*regexp.Regexp) {
	s.prompts[k] = regs
}

func (s *opFortinet) SetErrPatterns(regs []*regexp.Regexp) {
	s.errs = regs
}

func (s *opFortinet) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.transitions[k]; ok {
		return v
	}
	return nil
}

func (s *opFortinet) GetEncoding() string {
	return ""
}

func (s *opFortinet) GetExcludes() []*regexp.Regexp {
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

func (s *opFortinet) registerTransition(src, dst string) {
	k := src + "->" + dst

	if src == dst {
		// do nothing
		s.transitions[k] = []string{}
		return
	}

	// vdom -> login
	// global -> login
	if dst == "login" {
		// dst is login
		// just end from current mode
		s.transitions[k] = []string{"end"}
		return
	}

	// login -> global
	if src == "login" && dst == "global" {
		s.transitions[k] = []string{"config global"}
		return
	}

	// login -> vdom
	if src == "login" { // dst not login and not global
		// login to vdom
		s.transitions[k] = []string{"config vdom\n\t" +
			"edit " + dst +
			``}
		return
	}

	// vdom -> global == vdom -> login -> global
	if dst == "global" {
		s.transitions[k] = []string{"end\nconfig global"}
		return
	}

	// global -> vdom == global -> login -> vdom
	// vdomA -> vdomB == vdomA -> login -> vdomB
	s.transitions[k] = []string{"end\nconfig vdom\n\t" +
		"edit " + dst +
		``}
	return
}

func (s *opFortinet) RegisterMode(req *protocol.CliRequest) error {
	if s.GetPrompts(req.Mode) != nil {
		return nil
	}
	// no pattern for this mode
	// try insert
	logs.Info(req.LogPrefix, "registering pattern for mode", req.Mode)
	s.prompts[req.Mode] = []*regexp.Regexp{
		regexp.MustCompile(`[[:alnum:]]{1,}[[:alnum:]-_]{0,} \(` + req.Mode + `\) # $`),
	}
	// register transtions
	// someelse vdom/global mode may have been registered, but no transition made
	for k := range s.prompts {
		s.registerTransition(k, req.Mode)
		s.registerTransition(req.Mode, k)
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
		modes := ssh.TerminalModes{
			ssh.ECHO: 1, // enable echoing
		}
		if err := session.RequestPty("vt100", 0, 0, modes); err != nil {
			return nil, nil, nil, fmt.Errorf("request pty failed, %s", err)
		}
		if err := session.Shell(); err != nil {
			session.Close()
			return nil, nil, nil, fmt.Errorf("create shell failed, %s", err)
		}
		return r, w, session, nil
	}
}
