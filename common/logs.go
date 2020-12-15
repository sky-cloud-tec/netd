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

// LogConfig is the config struct for rrframework/logs
type LogConfig struct {
	Adaptor  string `json:"adaptor"`
	Filepath string `json:"filepath"`
	Level    string `json:"level"`
	MaxSize  int    `json:"maxsize"`
}

// AppConfig contains app config items
type AppConfig struct {
	Confidence int `json:"confidence"`
}

// AppConfigInstance ...
var AppConfigInstance *AppConfig
