## NetD

[![Build Status](https://travis-ci.org/sky-cloud-tec/netd.svg?branch=master)](https://travis-ci.org/sky-cloud-tec/netd)

#### how to run
```
go build .
./netd  --loglevel DEBUG jrpc
```

#### Usages
```go
	client, err := net.Dial("tcp", "localhost:8188")
	// Synchronous call
	args := &protocol.CliRequest{
		Device:  "juniper-test",
		Vendor:  "juniper",
		Type:    "srx",
		Version: "6.0",
		Address: "192.168.1.252:22",
		Auth: protocol.Auth{
			Username: "xx",
			Password: "xx",
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

#### Cli modes
* juniper
    * srx
        * login
        * configure
        * configure_private
        * configure_exclusive
* cisco

#### device support list
* juniper
    * srx
    * ssg
* cisco
    * asa
    * ios switch
    * nx-os switch
* fortinet
    * fortigate
* paloalto
* huawei
    * usg
