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

package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"

	"github.com/sky-cloud-tec/netd/protocol"
)

func main() {
	client, err := net.Dial("tcp", "localhost:8188")
	// Synchronous call
	// args := &protocol.CliRequest{
	// 	Device:  "topsec-telnet-test",
	// 	Vendor:  "topsec",
	// 	Type:    "NGFW4000",
	// 	Version: "xx",
	// 	Address: "10.20.20.200:22",
	// 	Auth: protocol.Auth{
	// 		Username: "superman",
	// 		Password: "Passw0rd.123",
	// 	},
	// 	Commands: []string{"show-running"},
	// 	Protocol: "telnet",
	// 	Mode:     "login",
	// 	Timeout:  30,
	// }
	args := &protocol.CliRequest{
		Device:  "cisco-telnet-test",
		Vendor:  "cisco",
		Type:    "asa",
		Version: "9.6",
		Address: "192.168.1.243:23",
		Auth: protocol.Auth{
			Username: "admin",
			Password: "abc@123",
		},
		Commands: []string{"show running-config"},
		Protocol: "telnet",
		Mode:     "login",
		Timeout:  30,
	}
	var reply protocol.CliResponse
	c := jsonrpc.NewClient(client)
	err = c.Call("CliHandler.Handle", args, &reply)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}
