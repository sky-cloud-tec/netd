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
	"io"
	"log"
	"regexp"

	"github.com/sky-cloud-tec/netd/protocol"
	"github.com/songtianyi/rrframework/logs"
	"golang.org/x/crypto/ssh"
)

// Operator interface
type Operator interface {
	GetTransitions(c, t string) []string
	GetPrompts(m string) []*regexp.Regexp
	GetErrPatterns() []*regexp.Regexp
	GetSSHInitializer() SSHInitializer
	GetLinebreak() string
	GetStartMode() string
}

var (
	// OperatorManagerInstance is OperatorManager instance
	OperatorManagerInstance *OperatorManager
)

// SSHInitializer ssh init func generator
type SSHInitializer func(*ssh.Client, *protocol.CliRequest) (io.Reader, io.WriteCloser, *ssh.Session, error)

func init() {
	OperatorManagerInstance = &OperatorManager{
		operatorMap: make(map[string]Operator, 0),
	}
}

// OperatorManager manager cli operators
type OperatorManager struct {
	operatorMap map[string]Operator // operatorMap mapping vendor.type.version to operator
}

// Get method return Operator instance by string
func (s *OperatorManager) Get(t string) Operator {
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
func (s *OperatorManager) Register(pattern string, o Operator) {
	logs.Info("Registering op", pattern, o)
	if _, ok := s.operatorMap[pattern]; ok {
		log.Fatal("pattern", pattern, "registered")
	}
	s.operatorMap[pattern] = o
}
