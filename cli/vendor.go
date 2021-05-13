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

package cli

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/sky-cloud-tec/netd/protocol"
	"github.com/songtianyi/rrframework/logs"
	"github.com/ziutek/telnet"
	"golang.org/x/crypto/ssh"
)

var (
	// VendorManagerInstance is VendorManager instance
	VendorManagerInstance *VendorManager
)

func init() {
	VendorManagerInstance = &VendorManager{
		operatorMap: make(map[string]*Vendor, 0),
	}
}

// VendorManager manager cli operators
type VendorManager struct {
	operatorMap map[string]*Vendor // operatorMap mapping vendor.type.version to operator
}

// Get method return Operator instance by string
func (s *VendorManager) Get(t string) *Vendor {
	for k, v := range s.operatorMap {
		logs.Debug("[ matching ]", k, t)
		if regexp.MustCompile(k).MatchString(t) {
			logs.Debug("[ matched ]", k, t)
			return v
		}
	}
	return nil
}

// Register do operator registration
func (s *VendorManager) Register(pattern string, o *Vendor) {
	logs.Info("Registering op", pattern, o)
	if _, ok := s.operatorMap[pattern]; ok {
		log.Fatal("pattern", pattern, "registered")
	}
	s.operatorMap[pattern] = o
}
// Vendor define the data type to be transferred
type Vendor struct {
	LineBreak    string // /r/n \n
	Transitions  map[string][]string
	Prompts      map[string][]*regexp.Regexp
	Excludes     []*regexp.Regexp
	Errs         []*regexp.Regexp
	Encoding     string
	Confidence   int
	StartMode    string
	CfgDebugFlag int
	CfgDebugDir  string
	Echo         bool
}

// SSHInitializer ssh session init func
type SSHInitializer func(*ssh.Client, *protocol.CliRequest) (io.Reader, io.WriteCloser, *ssh.Session, error)

// TELNETInitializer telnet conn init func
type TELNETInitializer func(*telnet.Conn, *protocol.CliRequest) error

// GetPrompts 判断Prompts是否有值存在，存在初始化变量v
func (s *Vendor) GetPrompts(k string) []*regexp.Regexp {
	if v, ok := s.Prompts[k]; ok {
		return v
	}
	return nil
}
// SetPrompts 初始化Prompts值
func (s *Vendor) SetPrompts(k string, regs []*regexp.Regexp) {
	s.Prompts[k] = regs
}
// SetErrPatterns 初始化Errs值
func (s *Vendor) SetErrPatterns(regs []*regexp.Regexp) {
	s.Errs = regs
}
// GetTransitions 判断Transitions是否有值存在，存在初始化变量v
func (s *Vendor) GetTransitions(c, t string) []string {
	k := c + "->" + t
	if v, ok := s.Transitions[k]; ok {
		return v
	}
	return nil
}
// GetEncoding 返回Encoding实例
func (s *Vendor) GetEncoding() string {
	return s.Encoding
}
// GetExcludes 返回Excludes实例
func (s *Vendor) GetExcludes() []*regexp.Regexp {
	return s.Excludes
}
// GetErrPatterns 返回Errs对应的正则表达
func (s *Vendor) GetErrPatterns() []*regexp.Regexp {
	return s.Errs
}
// GetStartMode 返回StartMode实例
func (s *Vendor) GetStartMode() string {
	return s.StartMode
}
// GetLinebreak 返回Linebreak实例
func (s *Vendor) GetLinebreak() string {
	return s.LineBreak
}

func (s *Vendor) registerTransition(src, dst string) {
	k := src + "->" + dst

	if src == dst {
		// do nothing
		s.Transitions[k] = []string{}
		return
	}

	// vdom -> login
	// global -> login
	if dst == "login" {
		// dst is login
		// just end from current mode
		s.Transitions[k] = []string{"end"}
		return
	}

	// login -> global
	if src == "login" && dst == "global" {
		s.Transitions[k] = []string{"config global"}
		return
	}

	// login -> vdom
	if src == "login" { // dst not login and not global
		// login to vdom
		s.Transitions[k] = []string{"config vdom\n\t" +
			"edit " + dst +
			``}
		return
	}

	// vdom -> global == vdom -> login -> global
	if dst == "global" {
		s.Transitions[k] = []string{"end\nconfig global"}
		return
	}

	// global -> vdom == global -> login -> vdom
	// vdomA -> vdomB == vdomA -> login -> vdomB
	s.Transitions[k] = []string{"end\nconfig vdom\n\t" +
		"edit " + dst +
		``}
	return
}

// RegisterMode 注册模式
func (s *Vendor) RegisterMode(req *protocol.CliRequest) error {
	if strings.ToLower(req.Vendor) != "fortinet" {
		return nil
	}
	if s.GetPrompts(req.Mode) != nil {
		return nil
	}
	// no pattern for this mode
	// try insert
	logs.Info(req.LogPrefix, "registering pattern for mode", req.Mode)
	s.Prompts[req.Mode] = []*regexp.Regexp{
		regexp.MustCompile(`[[:alnum:]]{1,}[[:alnum:]-_]{0,} \(` + req.Mode + `\) (#|\$) $`),
	}
	// register transtions
	// someelse vdom/global mode may have been registered, but no transition made
	for k := range s.Prompts {
		s.registerTransition(k, req.Mode)
		s.registerTransition(req.Mode, k)
	}
	logs.Debug(req.LogPrefix, s)
	return nil
}

// GetSSHInitializer 获取ssh连接通道
func (s *Vendor) GetSSHInitializer() SSHInitializer {
	return func(c *ssh.Client, req *protocol.CliRequest) (io.Reader, io.WriteCloser, *ssh.Session, error) {

		// consider vdom
		if strings.ToLower(req.Vendor) == "fortinet" {
			if err := s.RegisterMode(req); err != nil {
				return nil, nil, nil, err
			}
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
		if s.Echo {
			modes := ssh.TerminalModes{
				ssh.ECHO: 1, // enable echoing
			}
			if err := session.RequestPty("vt100", 0, 2000, modes); err != nil {
				return nil, nil, nil, fmt.Errorf("request pty failed, %s", err)
			}
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
