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

package protocol

import "time"

// CliRequest structure
type CliRequest struct {
	Vendor    string        `json:"vendor"`    // device vendor
	Type      string        `json:"type"`      // device type
	Version   string        `json:"version"`   // device os version
	Device    string        `json:"device"`    // device identity, uuid, hostname, etc.
	Mode      string        `json:"mode"`      // target mode
	Protocol  string        `json:"protocol"`  // telnet or ssh
	Auth      Auth          `json:"auth"`      // username and password
	Address   string        `json:"address"`   // host:port eg. 192.168.1.101:22
	Commands  []string      `json:"commands"`  // cli commands
	Format    string        `json:"format"`    //req format like xml,set
	Timeout   time.Duration `json:"timeout"`   // req timeout setting
	LogPrefix string        `json:"logPrefix"` // log prefix
	EnablePwd string        `json:"enablePwd"` // enable password for cisco devices
	Session   string        `json:"session"`   // session uuid
}

// Auth struct
type Auth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// CliResponse ...
type CliResponse struct {
	Retcode int
	Message string
	Device  string
	CmdsStd map[string]string
}
