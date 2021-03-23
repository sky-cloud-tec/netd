## NetD

[![Build Status](https://travis-ci.org/sky-cloud-tec/netd.svg?branch=master)](https://travis-ci.org/sky-cloud-tec/netd)

NetD 是深圳天元云科技开源的一款针对网络运维领域的命令执行器, 已成功对接数十种不同品牌版本的防火墙。

#### 特性
- 基于配置，无需修改和重新编译源码即可支持一种新的品牌或版本的硬件
- 支持 ssh, telnet
- 支持`模式`, 可以在指定模式下执行指定命令, 比如在普通模式下执行 `show running-config`, 在 `configure` 模式下执行 `set xxx` 等命令
- 自动切换模式, 无需手动执行 `configure`, `exit` 等模式切换命令
- 自动翻页, 当出现 `--More--` 时, 自动进行翻页获取完整配置
- 提供 jrpc 接口

#### 如何运行
加载默认配置
```
go build .
./netd
```

加载指定配置
```
go build .
./netd --cfg /path/to/cfg.ini
```

#### 如何添加一个新的设备支持

NetD 使用 `ini` 格式的配置文本, 以 Cisco 为例, 需要添加以下内容:
``` ini
; --- asa ---
; (Required) 添加一个新的 section, section name 为匹配设备的厂商.型号.版本的正则表达式
[(?i)cisco\.asa[a-z]{0,}\.(9|[0-9]{1,})\..*]
; linebreak, default is unix
; (Optional) 指定换行符
linebreak = windows
; prompts
; (Required) 提供用于匹配 prompt 的正则表达式, prompt. 的前缀是固定的，后面的名称可以自定义
prompt.login = "[[:alnum:]]{1,}(-[[:alnum:]]+){0,}> $"
prompt.login_enable = [[:alnum:]]{1,}(-[[:alnum:]]+){0,}# $
prompt.configure_terminal = "[[:alnum:]]{1,}(-[[:alnum:]]+){0,}\(config\)# $"
; modes
; (Optonal) 指定登陆设备后的起始模式, 因为设备差异和用户配置上的差异，不同设备登陆之后的默认模式是不同的
start = "login"
; (Required) 指定各个模式所使用的 prompt 规则, mode. 前缀固定
mode.login_or_login_enable = prompt.login, prompt.login_enable
mode.login = prompt.login
mode.login_enable = prompt.login_enable
mode.configure_terminal = prompt.configure_terminal
; transtions
; (Optional) 如果有多个模式，则必须指定, 提供模式之间切换时所执行的命令
transition.login_enable.configure_terminal = "configure terminal"
transition.configure_terminal.login_enable = "exit"
; error pattern
; (Required) 提供匹配命令执行错误的匹配规则
errs = "^ERROR: .*$"
; predefined encoding
; (Optional) 提供防火墙输出的编码格式
encoding = ""
; cancel more
; (Optional) 提供取消分页的命令及执行此命令所在的模式
cancel.login = "terminal pager 0", "terminal pager lines 0"
; debugging
; (Optional) 是否将命令的输出写入到文件
debug.cfg = false
; (Optional) 命令输出写入到文件的存放目录
debug.cfg_dir = /var/log/netd/cfgs
```

#### 如何调用
```go
    // go jrpc example
	client, err := net.Dial("tcp", "localhost:8188")
	// Synchronous call
	args := &protocol.CliRequest{
		Device:  "juniper-test",
		Vendor:  "juniper",
		Type:    "srx",
		Version: "6.0",
		Address: "192.168.1.252:22",
		Auth: protocol.Auth{
			Username: "xxx",
			Password: "xxx",
		},
		Commands: []string{"set security address-book global address WS-100.2.2.46_32 wildcard-address 100.2.2.46/255.255.255.255"},
		Protocol: "ssh",
		Mode:     "configure_private",
		Timeout:  30, // seconds
	}
	var reply protocol.CliResponse
	c := jsonrpc.NewClient(client)
	err = c.Call("CliHandler.Handle", args, &reply)
```
check [jrpc test](https://github.com/sky-cloud-tec/netd/blob/master/ingress/jrpc_test.go) file for more details

#### 常见问题
1. 命令执行错误，接口却没报错
> 错误输出没有被 NetD 捕捉到，可以通过在配置文件中添加 `errs` 的正则来解决

2. 命令执行超时
> NetD 是的基本原理是模式人的操作行为，人是通过 `prompt` 的出现来判断命令的执行是否结束，NetD 是一样的。将 log 的 level 设置为 `DEBUG`, 复现问题，如果观察到日志在打印完 prompt 之后没有其它输出，之后超时，则是因为 prompt 的匹配有问题，可以通过添加新的 prompt 或者修改已有的 prompt 来解决，最直接的解决方式是将整个 prompt 复制并放进配置里。

3. 命令不支持
> 不同型号版本的硬件所使用的命令差异较大，即使是同一品牌，同一型号的硬件，在不同版本里所使用的命令也有差异。直接在配置文件里修改原有命令，或者新增一个 section 来解决。

4. 修改配置后没生效
> 重启 NetD. NetD 不支持动态修改配置，需要手动重启完成配置的更新。注意，重启会导致正在执行的命令中断，所以传入 NetD 的命令最好是可重复执行的。