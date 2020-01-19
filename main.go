package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sky-cloud-tec/netd/common"
	"github.com/sky-cloud-tec/netd/ingress"

	"github.com/songtianyi/rrframework/logs"
	"github.com/urfave/cli"
)

// AppConfig app configurations
type AppConfig struct {
	logCfg *common.LogConfig
}

var appConfig *AppConfig

func init() {
	appConfig = &AppConfig{
		logCfg: &common.LogConfig{},
	}
}

func initLogger() error {
	// set logger
	property := `{"filename": "` + appConfig.logCfg.Filepath +
		`", "maxlines" : 10000000, "maxsize": ` + strconv.Itoa(appConfig.logCfg.MaxSize) + `}`
	fmt.Println(property, appConfig.logCfg)
	logs.SetLevel(common.MapStringToLevel[appConfig.logCfg.Level])
	return logs.SetLogger("file", property)

}

func jrpcHandler(c *cli.Context) error {
	// init logger
	if err := initLogger(); err != nil {
		return err
	}
	// init jrpc
	jrpc, _ := ingress.NewJrpc(c.String("addr"))
	jrpc.Register(new(ingress.CliHandler))
	if err := jrpc.Serve(); err != nil {
		return err
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = `NetD make network device operations easy!
	It's a dammon app which allow you to run cli commands through grpc, amqp etc.`
	app.Version = "2.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "songtianyi",
			Email: "songtianyi@sky-cloud.net",
		},
	}
	app.Copyright = "Copyright (c) 2017-2019 sky-cloud.net"
	app.Commands = []cli.Command{
		{
			Name:    "jrpc",
			Aliases: []string{"jrpc"},
			Usage:   "Run netd with jrpc ingress",
			Action:  jrpcHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address, addr",
					Value: "0.0.0.0:8188",
					Usage: "jprc listen address",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "logfile, lf",
			Value:       "/var/log/netd/netd.log",
			Usage:       "logfile path",
			Destination: &appConfig.logCfg.Filepath,
		},
		cli.StringFlag{
			Name:        "loglevel, ll",
			Value:       "INFO",
			Usage:       "log level, EMERGENCY|ALERT|CRITICAL|ERROR|WARNING|NOTICE|INFO|DEBUG",
			Destination: &appConfig.logCfg.Level,
		},
		cli.IntFlag{
			Name:        "maxsize, ms",
			Value:       10240000,
			Usage:       "log file max size",
			Destination: &appConfig.logCfg.MaxSize,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
