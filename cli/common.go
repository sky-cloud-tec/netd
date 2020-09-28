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
	"bytes"
	"errors"
	"github.com/songtianyi/rrframework/logs"
	"io"
	"regexp"
)

// AnyMatch return true if any regex in patterns matches the input string
func AnyMatch(patterns []*regexp.Regexp, s string) bool {
	if patterns == nil {
		return false
	}
	for _, v := range patterns {
		if v == nil {
			continue
		}
		matches := v.FindStringSubmatch(s)
		if len(matches) > 0 {
			return true
		}
	}
	return false
}

// Match return true if regexp p matches input string s
func Match(p *regexp.Regexp, s string) bool {
	if p == nil {
		return false
	}
	matches := p.FindStringSubmatch(s)
	return len(matches) > 0
}

// ReadStringUntil read string from reader until specified regex pattern matched
func ReadStringUntil(r io.Reader, p *regexp.Regexp) (string, error) {
	if p == nil || r == nil {
		return "", errors.New("p or r nil")
	}
	buf := make([]byte, 128)
	var (
		wbuf   bytes.Buffer
		errRes error
	)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			logs.Info("EOF")
			break
		} else if err != nil {
			// something wrong
			logs.Error("read error:", err)
			errRes = err
			break
		}
		// print received content
		logs.Debug("(", n, ")", string(buf[:n]))
		// write received content to whole document buffer
		wbuf.Write(buf[:n])
		if Match(p, wbuf.String()) {
			break
		}
	}
	return wbuf.String(), errRes
}

// ReadStringUntilError read string from reader until specified regex pattern or error matched
func ReadStringUntilError(r io.Reader, p *regexp.Regexp, e *regexp.Regexp) (string, error) {
	if p == nil || r == nil {
		return "", errors.New("p or r nil")
	}
	buf := make([]byte, 128)
	var (
		wbuf   bytes.Buffer
		errRes error
	)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			logs.Info("EOF")
			break
		} else if err != nil {
			// something wrong
			logs.Error("read error:", err)
			errRes = err
			break
		}
		// print received content
		logs.Debug("(", n, ")", string(buf[:n]))
		// write received content to whole document buffer
		wbuf.Write(buf[:n])
		if Match(p, wbuf.String()) {
			break
		}
		if Match(e, wbuf.String()) {
			errRes = errors.New("error matched")
			break
		}
	}
	return wbuf.String(), errRes

}
