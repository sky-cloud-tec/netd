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
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/ssh"
)

func TestConnDial(t *testing.T) {

	Convey("dial with ssh-dss", t, func() {
		sshConfig := &ssh.ClientConfig{
			User:            "admin",
			Auth:            []ssh.AuthMethod{ssh.Password("r00tme")},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}
		sshConfig.SetDefaults()
		sshConfig.Ciphers = append(sshConfig.Ciphers, []string{"aes128-cbc", "3des-cbc"}...)
		sshConfig.KeyExchanges = append(sshConfig.KeyExchanges, []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group-exchange-sha256"}...)
		client, err := ssh.Dial("tcp", "10.88.88.229:22", sshConfig)
		fmt.Println(err, client)
		So(
			err,
			ShouldBeNil,
		)
		So(
			client,
			ShouldNotBeNil,
		)
	})
}
