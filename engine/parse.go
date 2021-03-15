package engine

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/netd/common"
	"github.com/sky-cloud-tec/netd/ingress"
	"github.com/songtianyi/rrframework/logs"

	"gopkg.in/ini.v1"
)

func initLogger(cfg *ini.File) error {
	// set logger
	logSec := cfg.Section("log")
	property :=
		`{"filename": "` + logSec.Key("path").MustString("/var/log/netd/netd.log") +
			`", "maxlines" : 10000000, "maxsize": ` + strconv.Itoa(logSec.Key("max_size").MustInt(10240000)) + `}`
	fmt.Println(property, logSec)
	logs.SetLevel(common.MapStringToLevel[logSec.Key("level").MustString("INFO")])
	return logs.SetLogger("file", property)
}

func startJrpc(cfg *ini.File) error {
	// init jrpc
	ingressSec := cfg.Section("ingress")
	jrpc, _ := ingress.NewJrpc(ingressSec.Key("jrpc.addr").MustString("0.0.0.0:8188"))
	jrpc.Register(new(ingress.CliHandler))
	if err := jrpc.Serve(); err != nil {
		return err
	}
	return nil
}

func registerOp(cfg *ini.File) error {
	for _, sec := range cfg.Sections() {
		if sec.Name() == "log" || sec.Name() == "ingress" || sec.Name() == "app" {
			continue
		}

		for k := range sec.KeyStrings() {
			fmt.Println(k)
		}
		op := &cli.Vendor{
			Transitions: map[string][]string{},
			Prompts:     map[string][]*regexp.Regexp{},
			Errs:        []*regexp.Regexp{},
			LineBreak:   sec.Key("linebeak").MustString("\r\n"),
		}
		cli.VendorManagerInstance.Register(sec.Name(), op)
	}
	return nil
}

func LoadCfg(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	// init logger
	if err := initLogger(cfg); err != nil {
		return err
	}
	registerOp(cfg)
	startJrpc(cfg)
	return nil
}
