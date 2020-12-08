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

package common

import (
	"bytes"
	"github.com/songtianyi/rrframework/logs"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

var (
	// MapStringToLevel trans string to logs level enum
	MapStringToLevel = map[string]int{
		"EMERGENCY": logs.LevelEmergency,
		"ALERT":     logs.LevelAlert,
		"CRITICAL":  logs.LevelCritical,
		"ERROR":     logs.LevelError,
		"WARNING":   logs.LevelWarning,
		"NOTICE":    logs.LevelNotice,
		"INFO":      logs.LevelInformational,
		"DEBUG":     logs.LevelDebug,
	}
)

// GbkToUtf8 transfrom gbk byte to utf8 byte
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk transform utf8 byte to gbk byte
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

var enc = map[string]string{
	"GB-18030": "GB18030",
}

// ConvToUTF8 convert any encoding type byte to utf8 byte
func ConvToUTF8(src string, b []byte) ([]byte, error) {
	if src == "" || src == "UTF-8" || src == "utf-8" {
		return b, nil
	}
	if v, ok := enc[src]; ok {
		src = v
	}
	e, err := ianaindex.MIB.Encoding(src)
	if err != nil {
		return b, err
	}
	reader := transform.NewReader(bytes.NewReader(b), e.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return b, err
	}
	return d, nil
}
