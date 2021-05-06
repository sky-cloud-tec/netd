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
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"testing"

	"github.com/sky-cloud-tec/netd/common"
	"github.com/sky-cloud-tec/netd/protocol"

	. "github.com/smartystreets/goconvey/convey"
)
func TestPower6000_Show(t *testing.T) {
	//
	Convey("show power", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)

		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "power-show-test",
			Vendor:  "venustech",
			Type:    "Power_V6000-F1120",
			Version: "3.0",
			Address: "10.88.88.240:22",
			Auth: protocol.Auth{
				Username: "administrator",
				Password: "Admin@r00tme",
			},
			Commands: []string{"asdf"},
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
		// So(
		// 	reply.CmdsStd,
		// 	ShouldNotBeNil,
		// )
		// So(
		// 	len(reply.CmdsStd) == 1,
		// 	ShouldBeTrue,
		// )
	})
}
func TestTSOS_Show(t *testing.T) {
	//
	Convey("show tsos", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "tsos-show-test",
			Vendor:  "venustech",
			Type:    "tsos",
			Version: "6.0",
			Address: "10.88.88.8:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "Admin@r00tme",
			},
			Commands: []string{"asdf"},
			Protocol: "ssh",
			Mode:     "configure",
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
		// So(
		// 	reply.CmdsStd,
		// 	ShouldNotBeNil,
		// )
		// So(
		// 	len(reply.CmdsStd) == 1,
		// 	ShouldBeTrue,
		// )
	})
}

func TestJuniperSrx_Set(t *testing.T) {
	//
	Convey("set juniper srx cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "juniper-srx-set-test",
			Vendor:  "juniper",
			Type:    "srx",
			Version: "6.0",
			Address: "192.168.1.252:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{"set security address-book global address WS-100.2.2.46_32 wildcard-address 100.2.2.46/32", "commit"},
			Protocol: "ssh",
			Mode:     "configure_private",
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
			len(reply.CmdsStd) == 2,
			ShouldBeTrue,
		)
	})
}

//
func TestJuniperSrx_Show(t *testing.T) {
	//
	Convey("show juniper srx configuration", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "juniper-srx-show-test",
			Vendor:  "juniper",
			Type:    "srx",
			Version: "6.0",
			Address: "192.168.1.252:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{"show configuration | display set | no-more"},
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
			reply.CmdsStd,
			ShouldNotBeNil,
		)
		So(
			len(reply.CmdsStd) == 1,
			ShouldBeTrue,
		)
	})
}

//
func TestJuniperSsg_Set(t *testing.T) {
	//
	Convey("set juniper ssg cli commands", t, func() {

		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "juniper-ssg-set-test",
			Vendor:  "juniper",
			Type:    "SSG",
			Version: "ScreenOS(6.1.0)",
			Address: "192.168.1.229:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`set policy id 999 name liao from Trust to DMZ addr1111 xdf xxx permit
				set policy id 999
				set service TCP-4444
				exit`},
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

//
func TestJuniperSsg_show(t *testing.T) {
	//
	Convey("show juniper ssg cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "juniper-ssg-show-test",
			Vendor:  "juniper",
			Type:    "SSG",
			Version: "ScreenOS(6.1.0)",
			Address: "192.168.1.229:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{`get config`},
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

//
func TestCiscoAsa_Show(t *testing.T) {
	//
	Convey("show cisco asa configuration", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "cisco-asa-show-test",
			Vendor:  "cisco",
			Type:    "asa",
			Version: "9.6(x)",
			Address: "192.168.1.238:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands:  []string{"show running-config"},
			Protocol:  "ssh",
			Mode:      "login_enable",
			EnablePwd: "",
			Timeout:   30,
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
			reply.CmdsStd,
			ShouldNotBeNil,
		)
		fmt.Println(reply.CmdsStd)
		So(
			len(reply.CmdsStd) == 1,
			ShouldBeTrue,
		)
	})
}

func TestCiscoAsa_Connect(t *testing.T) {
	//
	Convey("test connect", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "cisco-asa-show-test",
			Vendor:  "cisco",
			Type:    "asa",
			Version: "9.6(x)",
			Address: "192.168.1.238:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands:  []string{},
			Protocol:  "ssh",
			Mode:      "login",
			EnablePwd: "",
			Timeout:   30,
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
	})
}

func TestCiscoAsa_Set(t *testing.T) {

	Convey("set cisco asa configuration", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "cisco-asa-set-test",
			Vendor:  "cisco",
			Type:    "asa",
			Version: "9.6(x)",
			Address: "192.168.1.238:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{"object network cisco-asa-set-test\n  host 1.1.1.1\nexit", "no object network cisco-asa-set-test"},
			Protocol: "ssh",
			Mode:     "configure_terminal",
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
			reply.CmdsStd,
			ShouldNotBeNil,
		)
		So(
			len(reply.CmdsStd) == 2,
			ShouldBeTrue,
		)
	})
}

func TestPaloalto_Set(t *testing.T) {
	//
	Convey("set Paloalto cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "paloalto-set-test",
			Vendor:  "paloalto",
			Type:    "pan-os",
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

//
func TestPaloalto_Show(t *testing.T) {
	//
	Convey("show Paloalto cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "paloalto-show-test",
			Vendor:  "paloalto",
			Type:    "pan-os",
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

func TestPanorama_Show(t *testing.T) {
	//
	Convey("show Paloalto cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "Panorama-connect-test",
			Vendor:  "paloalto",
			Type:    "Panorama",
			Version: "8.1",
			Address: "192.168.1.167:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "admin",
			},
			Commands: []string{},
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
	})
}

//
func TestUSG6000V2_Set(t *testing.T) {
	//
	Convey("set USG6000V2 cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "usg6000v2-set-test",
			Vendor:  "huawei",
			Type:    "usg",
			Version: "V500R005C10",
			Address: "192.168.1.205:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "Admin@r00tme",
			},
			Commands: []string{
				`security-policy
				  rule name policy_test
				    destination-address 10.1.1.1 24
				    destination-address 3000::1 32
				    destination-address geo-location BeiJing
				    quit
				  quit`,
			},
			Protocol: "ssh",
			Mode:     "system_View",
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

func TestUSG_Set2(t *testing.T) {
	//
	Convey("set USG6000 cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "usg6000-set-test",
			Vendor:  "huawei",
			Type:    "usg",
			Version: "V100R001C30SPCa00",
			Address: "192.168.1.236:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "admin",
			},
			Commands: []string{
				`ip address-set demo_10.1.1.3 type object
				    address 10.1.1.2 255.255.255.0
				    quit`,
			},
			Protocol: "ssh",
			Mode:     "system_View",
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

//
func TestUSG6000V2_Show(t *testing.T) {
	//
	Convey("Show USG6000V2 cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "usg6000v2-show-test",
			Vendor:  "huawei",
			Type:    "usg",
			Version: "usg6000v2",
			Address: "192.168.1.205:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "Admin@r00tme",
			},
			Commands: []string{
				//`display security-policy rule policy_test`,
				`display current-configuration`,
			},
			Protocol: "ssh",
			Mode:     "system_View",
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
		//fmt.Println(reply.CmdsStd)
		So(
			len(reply.CmdsStd) == 1,
			ShouldBeTrue,
		)
	})
}

func TestUSG5500_Telnet_Show_OK(t *testing.T) {
	//
	Convey("Show USG6000V2 cli commands by telnet", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "usg5500-show-test",
			Vendor:  "huawei",
			Type:    "usg",
			Version: "usg5500",
			Address: "192.168.1.177:23",
			Auth: protocol.Auth{
				Username: "hwtel",
				Password: "lablab@123",
			},
			Commands: []string{
				`display current-configuration`,
			},
			Protocol: "telnet",
			Mode:     "system_View",
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
		//fmt.Println(reply.CmdsStd)
		So(
			len(reply.CmdsStd) == 1,
			ShouldBeTrue,
		)
	})
}

func TestUSG5500_Telnet_Show_NOT_OK(t *testing.T) {
	//
	Convey("Show USG6000V2 cli commands by telnet", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "usg5500-show-test",
			Vendor:  "huawei",
			Type:    "usg",
			Version: "usg5500",
			Address: "192.168.1.177:23",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "lablab@123",
			},
			Commands: []string{
				`display current-configuration`,
			},
			Protocol: "telnet",
			Mode:     "system_View",
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
		//fmt.Println(reply.CmdsStd)
		So(
			len(reply.CmdsStd) == 1,
			ShouldBeTrue,
		)
	})
}

//
func TestIos_Show(t *testing.T) {
	//
	Convey("show ios cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "ios-show-test",
			Vendor:  "cisco",
			Type:    "IOS",
			Version: "I86BI_LINUXL2-IPBASEK9-M",
			Address: "192.168.1.244:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`show version`,
			},
			Protocol:  "ssh",
			Mode:      "login_enable",
			EnablePwd: "r00tme",
			Timeout:   30,
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

//
func TestIos_Set(t *testing.T) {
	//
	Convey("set ios cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "ios-set-test",
			Vendor:  "cisco",
			Type:    "IOS",
			Version: "I86BI_LINUXL2-IPBASEK9-M",
			Address: "192.168.1.244:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`hostname switch`,
			},
			Protocol: "ssh",
			Mode:     "configure_terminal",
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

//
func TestNxos_Show(t *testing.T) {
	//
	Convey("show nxos cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "nxos-show-test",
			Vendor:  "cisco",
			Type:    "NX-OS",
			Version: "7.0(3)I5(2)",
			Address: "192.168.1.248:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`show aaa accounting`,
			},
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

//
func TestNxos_Set(t *testing.T) {
	//
	Convey("set nxos cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "nxos-set-test",
			Vendor:  "cisco",
			Type:    "NX-OS",
			Version: "7.0(3)I5(2)",
			Address: "192.168.1.248:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`hostname ab`,
			},
			Protocol: "ssh",
			Mode:     "configure_terminal",
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

func TestHillstone_show(t *testing.T) {

	Convey("show hillstone cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "hillstone-show-test",
			Vendor:  "hillstone",
			Type:    "SG-6000-VM01",
			Version: "5.5",
			Address: "192.168.1.199:22",
			Auth: protocol.Auth{
				Username: "hillstone",
				Password: "hillstone",
			},
			Commands: []string{
				`show version`,
			},
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

func TestHillstone_set(t *testing.T) {

	Convey("set hillstone cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "hillstone-set-test",
			Vendor:  "hillstone",
			Type:    "SG-6000-VM01",
			Version: "5.5",
			Address: "192.168.1.232:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`hostname hillstone`,
			},
			Protocol: "ssh",
			Mode:     "configure",
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

func TestFortinet_set(t *testing.T) {

	Convey("set fortinet cli commands in global domain", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "fortinet-set-test",
			Vendor:  "fortinet",
			Type:    "FortiGate-VM64-KVM",
			Version: "v5.6.x",
			Address: "192.168.1.237:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`set hostname fortinet239`,
			},
			Protocol: "ssh",
			Mode:     "global",
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

func TestFortinet_exceed_policy_num_err(t *testing.T) {

	Convey("set policy id > limit", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "fortinet-policy-max-id-test",
			Vendor:  "fortinet",
			Type:    "FortiGate-VM64-KVM",
			Version: "v5.6.x",
			Address: "192.168.1.239:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`config firewall policy
				edit "501"
					set name "port1_2_trust_49c3a"
					set srcintf "port1"
					set dstintf "trust"
					set srcaddr "Net-192.168.1.34_32"
					set dstaddr "WS-172.16.10.34_32"
					set service "TCP-2154"
					set action accept
					set schedule "always"
					set comments "create by NAP 25e107b0-6334-4cc1-9d2e-b1b0bafeab2a"
				next
			end`,
			},
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
			ShouldBeFalse,
		)
	})
}

func TestFortinet_show(t *testing.T) {

	Convey("show fortinet cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "fortinet-show-test",
			Vendor:  "fortinet",
			Type:    "FortiGate-VM64-KVM",
			Version: "v5.6.x",
			Address: "192.168.1.239:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`show full-configuration`,
			},
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

func TestFortinet_show_telnet(t *testing.T) {

	Convey("show fortinet cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "fortinet-show-test",
			Vendor:  "fortinet",
			Type:    "FortiGate-VM64-KVM",
			Version: "v5.6.x",
			Address: "192.168.1.239:23",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`show full-configuration`,
			},
			Protocol: "telnet",
			Mode:     "login",
			Timeout:  90,
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

func TestTopSec_show(t *testing.T) {

	Convey("show topsec cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "topsec-show-test",
			Vendor:  "topsec",
			Type:    "NGFW4000",
			Version: "TG5030",
			Address: "192.168.1.208:22",
			Auth: protocol.Auth{
				Username: "superman",
				Password: "multisync88",
			},
			Commands: []string{
				`show`,
			},
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
		fmt.Print(reply.CmdsStd)
	})
}

func TestTianyuanJuniper_show(t *testing.T) {

	Convey("show tianyuan juniper conf", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "tianyuan-juniper-show-test",
			Vendor:  "juniper",
			Type:    "srx",
			Version: "6.0",
			Address: "192.168.2.102:22",
			Auth: protocol.Auth{
				Username: "test",
				Password: "test123",
			},
			Commands: []string{
				`show configuration | no-more`,
			},
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
		fmt.Print(reply.CmdsStd)
	})
}

func TestZihuJuniper_show(t *testing.T) {

	Convey("show tianyuan juniper conf", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "zihu-show-test",
			Vendor:  "juniper",
			Type:    "srx",
			Version: "6.0",
			Address: "192.168.1.242:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`show configuration | no-more`,
			},
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
		fmt.Print(reply.CmdsStd)
	})
}

func TestJuniper_show_telnet(t *testing.T) {

	Convey("show juniper conf with telnet", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "juniper-show-test-telnet",
			Vendor:  "juniper",
			Type:    "srx",
			Version: "6.0",
			Address: "192.168.1.161",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`show configuration | no-more`,
			},
			Protocol: "telnet",
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
		fmt.Print(reply.CmdsStd)
	})
}

//
func TestSecPathFW2000_Set(t *testing.T) {
	//
	Convey("set SecPathFW2000 cli commands", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "SecPathFW2000-set-test",
			Vendor:  "h3c",
			Type:    "secpath",
			Version: "FW2000",
			Address: "192.168.1.203:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`object-group ip address Host_172.16.203.2
                 network host address 172.16.203.2
				 quit`,
			},
			Protocol: "ssh",
			Mode:     "system_View",
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

//
func TestSecPathFW2000_Show(t *testing.T) {
	//
	Convey("show SecPathFW2000 config", t, func() {
		client, err := net.Dial("tcp", "localhost:8188")
		So(
			err,
			ShouldBeNil,
		)
		// Synchronous call
		args := &protocol.CliRequest{
			Device:  "SecPathFW2000-set-test",
			Vendor:  "H3C",
			Type:    "SecPath",
			Version: "FW2000",
			Address: "192.168.1.203:22",
			Auth: protocol.Auth{
				Username: "admin",
				Password: "r00tme",
			},
			Commands: []string{
				`display current-configuration `,
			},
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
