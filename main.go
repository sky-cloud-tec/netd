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
	"log"
	"os"
	"time"

	"github.com/sky-cloud-tec/netd/engine"

	"github.com/urfave/cli"
)

var dst string

func main() {
	app := cli.NewApp()
	app.Usage = `NetD make network device operations easy!
	It's a dammon app which allow you to run cli commands through grpc, amqp(not support yet) etc.`
	app.Version = "2.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "songtianyi",
			Email: "songtianyi@sky-cloud.net",
		},
	}
	app.Copyright = "Copyright (c) 2017-2021 sky-cloud.net"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "cfg-path, cfg",
			Value:       "cfg.ini",
			Usage:       "configuration file path",
			Destination: &dst,
		},
	}
	app.Action = func(c *cli.Context) error {
		if err := engine.LoadCfg(dst); err != nil {
			return err
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
