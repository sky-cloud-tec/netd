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

package ingress

import (
	"net"
	"net/rpc/jsonrpc"
	"testing"

	"github.com/sky-cloud-tec/netd/common"
	"github.com/sky-cloud-tec/netd/protocol"

	. "github.com/smartystreets/goconvey/convey"
)

// func TestJuniperSrx_Set(t *testing.T) {

// 	Convey("set juniper srx cli commands", t, func() {
// 		client, err := net.Dial("tcp", "localhost:8088")
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		// Synchronous call
// 		args := &protocol.CliRequest{
// 			Device:  "juniper-srx-set-test",
// 			Vendor:  "juniper",
// 			Type:    "srx",
// 			Version: "6.0",
// 			Address: "192.168.1.252:22",
// 			Auth: protocol.Auth{
// 				Username: "admin",
// 				Password: "r00tme",
// 			},
// 			Commands: []string{"set security address-book global address WS-100.2.2.46_32 wildcard-address 100.2.2.46/32", "commit"},
// 			Protocol: "ssh",
// 			Mode:     "configure_private",
// 			Timeout:  30,
// 		}
// 		var reply protocol.CliResponse
// 		c := jsonrpc.NewClient(client)
// 		err = c.Call("CliHandler.Handle", args, &reply)
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		So(
// 			reply.Retcode == common.OK,
// 			ShouldBeTrue,
// 		)
// 		So(
// 			len(reply.CmdsStd) == 2,
// 			ShouldBeTrue,
// 		)
// 	})
// }

// func TestJuniperSrx_Show(t *testing.T) {

// 	Convey("show juniper srx configuration", t, func() {
// 		client, err := net.Dial("tcp", "localhost:8088")
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		// Synchronous call
// 		args := &protocol.CliRequest{
// 			Device:  "juniper-srx-show-test",
// 			Vendor:  "juniper",
// 			Type:    "srx",
// 			Version: "6.0",
// 			Address: "192.168.1.252:22",
// 			Auth: protocol.Auth{
// 				Username: "admin",
// 				Password: "r00tme",
// 			},
// 			Commands: []string{"show configuration | display set | no-more"},
// 			Protocol: "ssh",
// 			Mode:     "login",
// 			Timeout:  30,
// 		}
// 		var reply protocol.CliResponse
// 		c := jsonrpc.NewClient(client)
// 		err = c.Call("CliHandler.Handle", args, &reply)
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		So(
// 			reply.Retcode == common.OK,
// 			ShouldBeTrue,
// 		)
// 		So(
// 			reply.CmdsStd,
// 			ShouldNotBeNil,
// 		)
// 		So(
// 			len(reply.CmdsStd) == 1,
// 			ShouldBeTrue,
// 		)
// 	})
// }

// func TestJuniperSsg_Set(t *testing.T) {

// 	Convey("set juniper ssg cli commands", t, func() {
// 		client, err := net.Dial("tcp", "localhost:8088")
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		// Synchronous call
// 		args := &protocol.CliRequest{
// 			Device:  "juniper-ssg-set-test",
// 			Vendor:  "juniper",
// 			Type:    "SSG",
// 			Version: "ScreenOS(6.1.0)",
// 			Address: "192.168.1.229:22",
// 			Auth: protocol.Auth{
// 				Username: "admin",
// 				Password: "r00tme",
// 			},
// 			Commands: []string{
// 				`set policy id 999 name liao from Trust to DMZ addr1111 xdf xxx permit
// 				set policy id 999
// 				set service TCP-4444
// 				exit`},
// 			Protocol: "ssh",
// 			Mode:     "login",
// 			Timeout:  30,
// 		}
// 		var reply protocol.CliResponse
// 		c := jsonrpc.NewClient(client)
// 		err = c.Call("CliHandler.Handle", args, &reply)
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		So(
// 			reply.Retcode == common.OK,
// 			ShouldBeTrue,
// 		)
// 		So(
// 			len(reply.CmdsStd) == 1,
// 			ShouldBeTrue,
// 		)
// 	})
// }

// func TestJuniperSsg_show(t *testing.T) {

// 	Convey("show juniper ssg cli commands", t, func() {
// 		client, err := net.Dial("tcp", "localhost:8088")
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		// Synchronous call
// 		args := &protocol.CliRequest{
// 			Device:  "juniper-ssg-show-test",
// 			Vendor:  "juniper",
// 			Type:    "SSG",
// 			Version: "ScreenOS(6.1.0)",
// 			Address: "192.168.1.229:22",
// 			Auth: protocol.Auth{
// 				Username: "admin",
// 				Password: "r00tme",
// 			},
// 			Commands: []string{`get config`},
// 			Protocol: "ssh",
// 			Mode:     "login",
// 			Timeout:  30,
// 		}
// 		var reply protocol.CliResponse
// 		c := jsonrpc.NewClient(client)
// 		err = c.Call("CliHandler.Handle", args, &reply)
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		So(
// 			reply.Retcode == common.OK,
// 			ShouldBeTrue,
// 		)
// 		So(
// 			len(reply.CmdsStd) == 1,
// 			ShouldBeTrue,
// 		)
// 	})
// }

// func TestCiscoAsa_Show(t *testing.T) {

// 	Convey("show cisco asa configuration", t, func() {
// 		client, err := net.Dial("tcp", "localhost:8088")
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		// Synchronous call
// 		args := &protocol.CliRequest{
// 			Device:  "cisco-asa-show-test",
// 			Vendor:  "cisco",
// 			Type:    "asa",
// 			Version: "9.6(x)",
// 			Address: "192.168.1.238:22",
// 			Auth: protocol.Auth{
// 				Username: "admin",
// 				Password: "r00tme",
// 			},
// 			Commands: []string{"show running-config"},
// 			Protocol: "ssh",
// 			Mode:     "login_enable",
// 			Timeout:  30,
// 		}
// 		var reply protocol.CliResponse
// 		c := jsonrpc.NewClient(client)
// 		err = c.Call("CliHandler.Handle", args, &reply)
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		So(
// 			reply.Retcode == common.OK,
// 			ShouldBeTrue,
// 		)
// 		So(
// 			reply.CmdsStd,
// 			ShouldNotBeNil,
// 		)
// 		So(
// 			len(reply.CmdsStd) == 1,
// 			ShouldBeTrue,
// 		)
// 	})
// }

// func TestCiscoAsa_Set(t *testing.T) {

// 	Convey("set cisco asa configuration", t, func() {
// 		client, err := net.Dial("tcp", "localhost:8088")
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		// Synchronous call
// 		args := &protocol.CliRequest{
// 			Device:  "cisco-asa-set-test",
// 			Vendor:  "cisco",
// 			Type:    "asa",
// 			Version: "9.6(x)",
// 			Address: "192.168.1.238:22",
// 			Auth: protocol.Auth{
// 				Username: "admin",
// 				Password: "r00tme",
// 			},
// 			Commands: []string{"object network cisco-asa-set-test\n  host 1.1.1.1\nexit", "no object network cisco-asa-set-test"},
// 			Protocol: "ssh",
// 			Mode:     "configure_terminal",
// 			Timeout:  30,
// 		}
// 		var reply protocol.CliResponse
// 		c := jsonrpc.NewClient(client)
// 		err = c.Call("CliHandler.Handle", args, &reply)
// 		So(
// 			err,
// 			ShouldBeNil,
// 		)
// 		So(
// 			reply.Retcode == common.OK,
// 			ShouldBeTrue,
// 		)
// 		So(
// 			reply.CmdsStd,
// 			ShouldNotBeNil,
// 		)
// 		So(
// 			len(reply.CmdsStd) == 2,
// 			ShouldBeTrue,
// 		)
// 	})
// }

func TestPaloalto_Set(t *testing.T) {

	Convey("set Paloalto cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8088")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "paloalto-set-test",
			Vendor:  "Paloalto",
			Type:    "Pan-OS",
			Version: "8.1",
			Address: "192.168.1.231:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`set deviceconfig system hostname PA-VM-1`,
				`commit`},
			Protocol: "ssh",
			Mode:     "configure",
			Timeout:  120, // commit
		}
		var reply protocol.CliResponse
		c := jsonrpc.NewClient(client)
		err = c.Call("CliHandler.Handle", args, &reply)
		So(
			err,
			ShouldBeNil,
		)
		So(
			reply.Retcode == common.OK,
			ShouldBeTrue,
		)
		So(
			len(reply.CmdsStd) == 2,
			ShouldBeTrue,
		)
	})
}

func TestPaloalto_Show(t *testing.T) {

	Convey("show Paloalto cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8088")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "paloalto-show-test",
			Vendor:  "Paloalto",
			Type:    "Pan-OS",
			Version: "8.1",
			Address: "192.168.1.231:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`
				show config running`},
			Protocol: "ssh",
			Mode:     "login",
			Timeout:  30,
		}
		var reply protocol.CliResponse
		c := jsonrpc.NewClient(client)
		err = c.Call("CliHandler.Handle", args, &reply)
		So(
			err,
			ShouldBeNil,
		)
		So(
			reply.Retcode == common.OK,
			ShouldBeTrue,
		)
		So(
			len(reply.CmdsStd) == 1,
			ShouldBeTrue,
		)
	})
}
