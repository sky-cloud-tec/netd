package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	fmt.Println(property)
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
		fmt.Println("[", sec.Name(), "]")
		if sec.Name() == "DEFAULT" || sec.Name() == "default" || sec.Name() == "log" || sec.Name() == "ingress" {
			continue
		}

		prompts_vars := make(map[string]*regexp.Regexp, 0)
		mode_prompts := make(map[string][]*regexp.Regexp, 0)
		trans := make(map[string][]string, 0)
		errs := make([]*regexp.Regexp, 0)
		excludes := make([]*regexp.Regexp, 0)
		for _, k := range sec.Keys() {
			fmt.Println(k.Name(), "=", k.String())
			parts := strings.Split(k.Name(), ".")
			switch parts[0] {
			case "prompt":
				// eg. prompt.active_login
				prompts_vars[k.Name()] = regexp.MustCompile(k.Value())
			case "excludes":
				for _, pname := range k.Strings(",") {
					// append single prompt regex to mode
					v, ok := prompts_vars[pname]
					if !ok {
						return fmt.Errorf("prompt %s not defined in previous section", pname)
					}
					excludes = append(excludes, v)
				}
			case "mode":
				mode := parts[1]
				if mode_prompts[mode] == nil {
					mode_prompts[mode] = make([]*regexp.Regexp, 0)
				}
				for _, pname := range k.Strings(",") {
					// append single prompt regex to mode
					v, ok := prompts_vars[pname]
					if !ok {
						return fmt.Errorf("prompt %s not defined in previous section", pname)
					}
					mode_prompts[mode] = append(mode_prompts[mode], v)
				}
			case "transition":
				direction := parts[1] + "->" + parts[2]
				trans[direction] = k.Strings(",")
			case "errors", "errs":
				// may enclosed by quotes
				for _, e := range k.Strings(",") {
					if (e[0] == '"' || e[0] == '\'') && e[len(e)-1] == e[0] {
						e = e[1 : len(e)-1]
					}
					errs = append(errs, regexp.MustCompile(e))
				}
			case "linebreak":
			case "encoding":
			case "start":
			case "cancel":
			case "debug":
			case "init":
			default:
				return fmt.Errorf("unsupported config key %s", parts[0])
			}
		}
		fmt.Println(sec.Key("linebreak").MustString("\n"))
		op := &cli.Vendor{
			Transitions: trans,
			Prompts:     mode_prompts,
			Excludes:    excludes,
			Errs:        errs,
			LineBreak:   sec.Key("linebreak").MustString("\n"),
			Encoding:    sec.Key("encoding").MustString(""),
			StartMode:   sec.Key("start").MustString("login"),
		}
		fmt.Println(op)
		cli.VendorManagerInstance.Register(sec.Name(), op)
	}
	return nil
}

func LoadCfg(path string) error {
	opts := ini.LoadOptions{
		Insensitive:         true,
		IgnoreInlineComment: true,
	}
	cfg, err := ini.LoadSources(opts, path)
	if err != nil {
		return err
	}
	// init logger
	if err := initLogger(cfg); err != nil {
		return err
	}
	if err := registerOp(cfg); err != nil {
		return err
	}
	startJrpc(cfg)
	return nil
}
