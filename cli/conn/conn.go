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

package conn

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/protocol"
	"github.com/songtianyi/rrframework/logs"

	"github.com/sky-cloud-tec/netd/common"
	"golang.org/x/crypto/ssh"

	"github.com/ziutek/telnet"
)

var (
	conns map[string]*CliConn
	semas map[string]chan struct{}
)

func init() {
	conns = make(map[string]*CliConn, 0)
	semas = make(map[string]chan struct{}, 0)
}

// CliConn cli connection
type CliConn struct {
	t    int                  // connection type 0 = ssh, 1 = telnet
	mode string               // device cli mode
	req  *protocol.CliRequest // cli request
	op   cli.Operator         // cli operator

	conn   *telnet.Conn // telnet connection
	client *ssh.Client  // ssh client

	session *ssh.Session   // ssh session
	r       io.Reader      // ssh session stdout
	w       io.WriteCloser // ssh session stdin

	formatSet bool
}

// Acquire cli conn
func Acquire(req *protocol.CliRequest, op cli.Operator) (*CliConn, error) {
	// limit concurrency to 1
	// there only one req for one connection always
	logs.Info(req.LogPrefix, "Acquiring sema...")
	if semas[req.Address] == nil {
		semas[req.Address] = make(chan struct{}, 1)
	}
	// try
	semas[req.Address] <- struct{}{}
	logs.Info(req.LogPrefix, "sema acquired")
	// no matter what going on next, sema should be released once
	if req.Mode == "" {
		req.Mode = op.GetStartMode()
	}
	// if cli conn already created
	if v, ok := conns[req.Address]; ok {
		v.req = req
		v.op = op
		logs.Info(req.LogPrefix, "cli conn exist")
		return v, nil
	}
	c, err := newCliConn(req, op)
	if err != nil {
		// sema will be released in parent func
		return nil, err
	}
	conns[req.Address] = c
	return c, nil
}

// Release cli conn
func Release(req *protocol.CliRequest) {
	if len(semas[req.Address]) > 0 {
		logs.Info(req.LogPrefix, "Releasing sema")
		<-semas[req.Address]
	}
	logs.Info(req.LogPrefix, "sema released")
}

func newCliConn(req *protocol.CliRequest, op cli.Operator) (*CliConn, error) {
	logs.Info(req.LogPrefix, "creating cli conn...")
	if strings.ToLower(req.Protocol) == "ssh" {
		sshConfig := &ssh.ClientConfig{
			User:            req.Auth.Username,
			Auth:            []ssh.AuthMethod{ssh.Password(req.Auth.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}
		sshConfig.SetDefaults()
		sshConfig.Ciphers = append(sshConfig.Ciphers, []string{"aes128-cbc", "3des-cbc"}...)
		sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group-exchange-sha256"}...)
		client, err := ssh.Dial("tcp", req.Address, sshConfig)
		if err != nil {
			logs.Error(req.LogPrefix, "dial", req.Address, "error", err)
			return nil, fmt.Errorf("dial %s error, %s", req.Address, err)
		}
		c := &CliConn{t: common.SSHConn, client: client, req: req, op: op, mode: op.GetStartMode()}
		if err := c.init(); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	} else if strings.ToLower(req.Protocol) == "telnet" {
		conn, err := telnet.DialTimeout("tcp", req.Address, 5*time.Second)
		if err != nil {
			return nil, fmt.Errorf("dial %s error, %s", req.Address, err)
		}
		c := &CliConn{t: common.TELNETConn, conn: conn, req: req, op: op, mode: op.GetStartMode()}
		return c, nil
	}
	return nil, fmt.Errorf("protocol %s not support", req.Protocol)
}

func (s *CliConn) heartbeat() {
	go func() {
		tick := time.Tick(30 * time.Second)
		for {
			select {
			case <-tick:
				// try
				logs.Info(s.req.LogPrefix, "Acquiring heartbeat sema...")
				semas[s.req.Address] <- struct{}{}
				logs.Info(s.req.LogPrefix, "heartbeat sema acquired")
				if _, err := s.writeBuff(""); err != nil {
					logs.Critical(s.req.LogPrefix, "heartbeat error,", err)
					s.Close()
					Release(s.req)
					return
				}
				if _, _, err := s.readBuff(); err != nil {
					logs.Critical(s.req.LogPrefix, "heartbeat error,", err)
					s.Close()
					Release(s.req)
					return
				}
				// OK
				Release(s.req)
			}
		}
	}()
}

func (s *CliConn) init() error {
	if s.t == common.SSHConn {
		f := s.op.GetSSHInitializer()
		var err error
		s.r, s.w, s.session, err = f(s.client, s.req)
		if err != nil {
			return err
		}
	}
	// read login prompt
	_, prompt, err := s.readBuff()
	if err != nil {
		return fmt.Errorf("read after login failed, %s", err)
	}
	// enable cases
	if s.mode == "login_or_login_enable" {
		// not sure what mode it is
		// check prompt
		loginPrompts := s.op.GetPrompts("login")
		if cli.Match(loginPrompts, prompt) {
			// in login, not enabled
			// eventually, we'll exec cmds in privileged mode
			s.mode = "login"
			if s.mode != s.req.Mode {
				// login is not the target mode, need transition
				// enter privileged mode
				if _, err := s.writeBuff("enable\r" + s.req.EnablePwd); err != nil {
					return fmt.Errorf("enter privileged mode err, %s", err)
				}
				s.mode = "login_enable"
				if _, _, err := s.readBuff(); err != nil {
					s.mode = "login"
					return fmt.Errorf("readBuff after enable err, %s", err)
				}
				if err := s.closePage(true); err != nil {
					return err
				}
			} // login is what you want, no close page here
		} else {
			// already in privileged mode, close page
			if err := s.closePage(true); err != nil {
				return err
			}
		}
	} else { //
		// special devices
		if strings.EqualFold(s.req.Vendor, "fortinet") && strings.EqualFold(s.req.Type, "fortigate-VM64-KVM") {
			pts := s.op.GetPrompts(s.req.Mode)
			if pts == nil {
				return fmt.Errorf("mode %s not registered", s.req.Mode)
			}
			// close page
			if !strings.Contains(pts[0].String(), s.req.Mode) {
				//non vdom
				s.closePage(true)
			} else {
				// vdom
				logs.Debug(s.req.LogPrefix, "entering domain global...")
				if _, err := s.writeBuff("config global"); err != nil {
					return err
				}
				if err := s.closePage(false); err != nil {
					return err
				}
				logs.Debug(s.req.LogPrefix, "exiting vdom global...")
				if _, err := s.writeBuff("end"); err != nil {
					return err
				}
				if _, _, err := s.readBuff(); err != nil {
					return err
				}
			}
		} else {
			// for any other non special devices
			// include cisco asa non login_or_login_enable mode
			if err := s.closePage(true); err != nil {
				return err
			}
		}
	}
	s.heartbeat()
	return nil
}

func (s *CliConn) closePage(drain bool) error {
	if strings.EqualFold(s.req.Vendor, "cisco") && (strings.EqualFold(s.req.Type, "asa") || strings.EqualFold(s.req.Type, "asav")) {
		// login mode no close page
		if s.mode == "login" {
			return nil
		}
		// ===config or normal both ok===
		// set terminal pager
		if _, err := s.writeBuff("terminal pager 0"); err != nil {
			return err
		}
		if _, _, err1 := s.readBuff(); err1 != nil {
			return err1
		}
		// set page lines
		if _, err := s.writeBuff("terminal pager lines 0"); err != nil {
			return err
		}
	} else if strings.EqualFold(s.req.Vendor, "cisco") && strings.EqualFold(s.req.Type, "ios") {
		if s.mode == "login" {
			return nil
		}
		if _, err := s.writeBuff("terminal length 0"); err != nil {
			return err
		}
	} else if strings.EqualFold(s.req.Vendor, "Paloalto") && strings.EqualFold(s.req.Type, "PAN-OS") {
		// set pager
		if _, err := s.writeBuff("set cli pager off"); err != nil {
			return err
		}
	} else if strings.EqualFold(s.req.Vendor, "hillstone") && strings.EqualFold(s.req.Type, "SG-6000-VM01") {
		// set pager
		if _, err := s.writeBuff("terminal length 0"); err != nil {
			return err
		}
	} else if strings.EqualFold(s.req.Vendor, "fortinet") && strings.EqualFold(s.req.Type, "fortigate-VM64-KVM") {
		// set console
		if _, err := s.writeBuff("config system console\n\tset output standard\nend"); err != nil {
			return err
		}
	} else {
		// we did not write any commands to these devices
		// so no need to readBuff
		return nil
	}
	if !drain {
		return nil
	}
	if _, _, err := s.readBuff(); err != nil {
		return err
	}
	return nil
}

// Close cli conn
func (s *CliConn) Close() error {
	delete(conns, s.req.Address)
	if s.t == common.TELNETConn {
		if s.conn == nil {
			logs.Info(s.req.LogPrefix, "telnet conn nil when close")
			return nil
		}
		return s.conn.Close()
	}
	if s.session != nil {
		if err := s.session.Close(); err != nil {
			return err
		}
	} else {
		logs.Notice(s.req.LogPrefix, "ssh session nil when close")
	}
	if s.client == nil {
		logs.Notice(s.req.LogPrefix, "ssh conn nil when close")
		return nil
	}
	return s.client.Close()
}

func (s *CliConn) read(buff []byte) (int, error) {
	if s.t == common.SSHConn {
		return s.r.Read(buff)
	}
	return s.conn.Read(buff)
}

func (s *CliConn) write(b []byte) (int, error) {
	if s.t == common.SSHConn {
		return s.w.Write(b)
	}
	return s.conn.Write(b)
}

type readBuffOut struct {
	err    error
	ret    string
	prompt string
}

// AnyPatternMatches return matched string slice if any pattern fullfil
func (s *CliConn) anyPatternMatches(t string, patterns []*regexp.Regexp) []string {
	for _, v := range patterns {
		matches := v.FindStringSubmatch(t)
		if len(matches) != 0 {
			return matches
		}
	}
	return nil
}

func (s *CliConn) readLines() *readBuffOut {
	buf := make([]byte, 1000)
	var (
		waitingString, lastLine string
		errRes                  error
		wbuf                    bytes.Buffer
	)
outside:
	for {
		n, err := s.read(buf) //this reads the ssh/telnet terminal
		if err != nil {
			// something wrong
			logs.Error(s.req.LogPrefix, "io.Reader read error,", err)
			errRes = err
			break
		}

		// print received content
		logs.Debug(s.req.LogPrefix, "(", n, ")", string(buf[:n]))

		// write received content to whole document buffer
		wbuf.Write(buf[:n])
		// slice alias
		rbuf := wbuf.Bytes()

		// reverse traversal
		// traverse lastline
		var beginIdx int
		for i := wbuf.Len() - 1; i >= 0; i-- {
			if rbuf[i] == '\n' || rbuf[i] == '\r' {
				beginIdx = i
				break
			}
		}
		testee := string(rbuf[beginIdx:])
		// check prompt patterns
		if s.op.GetPrompts(s.mode) == nil {
			logs.Error(s.req.LogPrefix, "no patterns for mode", s.mode)
			errRes = fmt.Errorf("no patterns for mode %s", s.mode)
			break outside
		}
		// test
		matches := s.anyPatternMatches(testee, s.op.GetPrompts(s.mode))
		if len(matches) > 0 {
			// test pass
			logs.Info(s.req.LogPrefix, "prompt matched", s.mode, ":", matches)
			// ignore prompt and break
			if beginIdx == 0 {
				lastLine = testee
			} else {
				// newline not include
				lastLine = string(rbuf[beginIdx+1:])
				// \n not include but \r maybe include in windows linebreak
				waitingString = string(rbuf[:beginIdx])
			}
			// break the out loop
			break outside
		}
		// not match
		// check buf size is large enough
		if cap(buf) == n {
			// buf full, it proves that maybe there are lots of more content out there
			// enlarge buf
			buf = make([]byte, 2*n)
		}
	}
	return &readBuffOut{
		errRes,
		waitingString,
		lastLine,
	}
}

// return cmd output, prompt, error
func (s *CliConn) readBuff() (string, string, error) {
	// buffered chan
	ch := make(chan *readBuffOut, 1)

	go func() {
		ch <- s.readLines()
	}()

	select {
	case res := <-ch:
		if res.err == nil {
			scanner := bufio.NewScanner(strings.NewReader(res.ret))
			for scanner.Scan() {
				matches := s.anyPatternMatches(scanner.Text(), s.op.GetErrPatterns())
				if len(matches) > 0 {
					logs.Info(s.req.LogPrefix, "err pattern matched,", res.ret)
					return "", res.prompt, fmt.Errorf("err pattern matched, %s", res.ret)
				}
			}
		}
		return res.ret, res.prompt, res.err
	case <-time.After(s.req.Timeout):
		return "", "", fmt.Errorf("read stdout timeout after %q", s.req.Timeout)
	}
}

func (s *CliConn) writeBuff(cmd string) (int, error) {
	if len(cmd) > 0 && cmd[len(cmd)-1] == '\n' {
		return s.write([]byte(cmd))
	}
	return s.write([]byte(cmd + s.op.GetLinebreak()))
}

// Exec execute cli cmds
func (s *CliConn) Exec() (map[string]string, error) {
	if err := s.beforeExec(); err != nil {
		logs.Error(s.req.LogPrefix, "beforeTransition error", err)
		return nil, fmt.Errorf("beforeTransition error, %s", err)
	}
	// transit to target mode
	if s.req.Mode != s.mode {
		s.op.RegisterMode(s.req)
		cmds := s.op.GetTransitions(s.mode, s.req.Mode)
		// use target mode prompt
		logs.Info(s.req.LogPrefix, s.mode, "-->", s.req.Mode)
		// transition back when it fail
		mt := s.mode
		s.mode = s.req.Mode
		for _, v := range cmds {
			logs.Info(s.req.LogPrefix, "exec", "<", v, ">")
			if _, err := s.writeBuff(v); err != nil {
				logs.Error(s.req.LogPrefix, "write buff failed,", err)
				s.mode = mt
				return nil, fmt.Errorf("write buff failed, %s", err)
			}
			_, _, err := s.readBuff()
			if err != nil {
				s.mode = mt
				logs.Error(s.req.LogPrefix, "readBuff failed,", err)
				return nil, fmt.Errorf("readBuff failed, %s", err)
			}
		}
	}
	if err := s.beforeExec(); err != nil {
		logs.Error(s.req.LogPrefix, "beforeExec error", err)
		return nil, fmt.Errorf("beforeExec error, %s", err)
	}
	cmdstd := make(map[string]string, 0)
	// do execute cli commands
	for _, v := range s.req.Commands {
		logs.Info(s.req.LogPrefix, "exec", "<", v, ">", "in", s.mode, "mode")
		if _, err := s.writeBuff(v); err != nil {
			logs.Error(s.req.LogPrefix, "write buff failed,", err)
			return cmdstd, fmt.Errorf("write buff failed, %s", err)
		}
		ret, _, err := s.readBuff()
		if err != nil {
			logs.Error(s.req.LogPrefix, "readBuff failed,", err)
			return cmdstd, fmt.Errorf("readBuff failed, %s", err)
		}
		cmdstd[v] = ret
	}
	return cmdstd, nil
}

func (s *CliConn) beforeExec() error {
	if s.req.Format == "" || s.formatSet {
		return nil
	}
	if strings.EqualFold(s.req.Vendor, "Paloalto") && strings.EqualFold(s.req.Type, "PAN-OS") {
		// set format
		// only for pa device
		mode := s.mode
		if s.mode != "login" {
			// transition to login first
			if _, err := s.writeBuff("exit"); err != nil {
				return err
			}
			s.mode = "login"
			if _, _, err := s.readBuff(); err != nil {
				s.mode = mode
				return err
			}
		}
		// set format
		if _, err := s.writeBuff("set cli config-output-format " + s.req.Format); err != nil {
			return err
		}
		if _, _, err := s.readBuff(); err != nil {
			return err
		}
		if mode != "login" {
			// transition back
			if _, err := s.writeBuff("configure"); err != nil {
				return err
			}
			s.mode = mode
			if _, _, err := s.readBuff(); err != nil {
				s.mode = "login"
				return err
			}
		}
	}
	return nil
}

func (s *CliConn) beforeTransition() error {
	if s.req.Format == "" {
		return nil
	}
	// if current in login mode
	if s.mode == "login" &&
		strings.EqualFold(s.req.Vendor, "Paloalto") && strings.EqualFold(s.req.Type, "PAN-OS") {
		// set format
		if _, err := s.writeBuff("set cli config-output-format " + s.req.Format); err != nil {
			return err
		}
		if _, _, err := s.readBuff(); err != nil {
			return err
		}
		s.formatSet = true
	}
	return nil
}
