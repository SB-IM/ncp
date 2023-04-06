# Node control protocol

Edge gateway core communication service

[![Build Status](https://github.com/sb-im/ncp/workflows/ci/badge.svg)](https://github.com/sb-im/ncp/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/sb-im/ncp/branch/master/graph/badge.svg)](https://codecov.io/gh/sb-im/ncp)
[![Go Report Card](https://goreportcard.com/badge/github.com/sb-im/ncp)](https://goreportcard.com/report/github.com/sb-im/ncp)
[![Go Reference](https://pkg.go.dev/badge/github.com/sb-im/ncp.svg)](https://pkg.go.dev/github.com/sb-im/ncp)
[![GitHub release](https://img.shields.io/github/tag/sb-im/ncp.svg?label=release)](https://github.com/sb-im/ncp/releases)
[![license](https://img.shields.io/github/license/sb-im/ncp.svg?maxAge=2592000)](https://github.com/sb-im/ncp/blob/master/LICENSE)

[README](README.md) | [中文文档](README_zh.md)

## What is ncp

This is a communication module, and there is no specific use specified

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Compile && Run](#compile--run)
* [Architecture design](#architecture-design)
* [NcpIO](#ncpio)
  * [Rules](#rules)
  * [Debug](#debug)
  * [Name](#name)
* [IOs](#ios)
  * [API](#api)
  * [tcpc](#tcpc)
  * [tcps](#tcps)
  * [exec](#exec)
  * [jsonrpc2](#jsonrpc2)
  * [logger](#logger)
  * [mqtt](#mqtt)
* [mqttd](#mqttd)
  * [mqttd-rpc](#mqttd-rpc)
  * [mqttd-status](#mqttd-status)
  * [mqttd-network](#mqttd-network)
  * [mqttd-trans](#mqttd-trans)
  * [mqttd-buildin](#mqttd-buildin)
* [Example](#example)
* [TODO](#todo)

<!-- vim-markdown-toc -->

## Compile && Run

* make
* go >= 1.13

```sh
make

# run
cp conf/config-dist.yml config.yml
ncp -c config.yml
```

## Architecture design

![architecture](docs/ncp-architecture.svg)

## NcpIO

This is the core of the entire application and can also be used as a library

The model used here is an `IO` module and a message center, where each `IO` module receives a message and broadcasts it to the other `IO` modules

The `rules` are used to determine whether the message is acceptable or not, and it is negative when no rule is matched (the message is not accepted)

This module can be used as a library.
Messages can be sent and received directly using the `api` interface

```yaml
ncpio:
  - type: tcpc
    name: airctl
    params: "localhost:8980"
    i_rules:
      - regexp: '.*"method": ?"(webrtc|history|ncp)".*'
        invert: true
      - regexp: '.*"method".*'
    o_rules:
      - regexp: '.*'
  - type: api
    name: airctl
```

### Rules

![ncpio-rules](docs/ncp-ncpio.svg)

* `o_rules` output the io rule (after recv)
* `i_rules` input the io rule (before send)

1. `i_rules` & `o_rules` use regexp matched, because message should plain
2. There can be more than one rule, match them accordingly, and as long as one of them is met, the rule will be allows
3. The `invert` field controls whether the regular match result should be inverted
4. When no rule is matched, the message does disallows

### Debug

The complexity of using the regular expression part is hard to debug. That's why debug mode was added

```bash
./ncp -debug
```

Every messages by `IO` need two rules group

Three log: (RECV, BROADCAST, SEND)

1. `RECV`: `IO`  When received the message
2. by `o_rules`
3. `BROADCAST`: When broadcasting a message to all `IO`
4. by `i_rules`
5. `SEND`: `IO` When can send a message

### Name

There can be more than one `IO` module of each type. You can give each one a name

By default, a unique name is automatically generated, incrementing from `0`, `1`, `2`.

It is recommended not to set the `name` field manually, if you set the same `name` field,
it will be placed in the broadcast message and sent to you by `name`.

The same `name` will not receive the broadcast message with the same name.
**Please check if you intentionally do not receive the broadcast, or if you accidentally set it to be the same**

## IOs

Type                 | Description
-------------------- | -----------
[API](#api)          | build-in interface
[tcpc](#tcpc)        | TCP socket client
[tcps](#tcps)        | TCP socket server
[jsonrpc2](#jsonrpc2)| JSONRPC 2.0 simulation
[logger](#logger)    | Record message log
[mqtt](#mqtt)        | Connect Mqtt Broker


At the core are six `IO` modules, with different combinations defined by configuration files to achieve different functions

Why is there no `udp server` and `udp client`?

There was no need for this at the time of development, some messages need to be verified for reliability. `udp` is not suitable

### API

Is a special interface

It is used when ncpio is used as a library, and through this API is an IO module

```go
import(
	ncpio "sb.im/ncp/ncpio"
)

func main()
	ncpios := ncpio.NewNcpIOs([]ncpio.Config{
		{
			Type: "api",
			IRules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
			ORules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
		},
		{
			Type:   "jsonrpc2",
			Params: `{"result":"ok"}`,
			IRules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
			ORules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
		},
	})
	go ncpios.Run(ctx)

	ncpio.I <- []byte(`{"jsonrpc":"2.0","method":"test","id":"test.0"}`)
	fmt.Println(string(<-ncpio.O))
}
```

Note: Unit tests are tested using this interface

There can be at most one API interface

### tcpc

As a tcp client to connect to the server, the most common working mode

`tcpc`, `tcps` message splitting using `\n` as separator, will automatically `chomp` `\r\n` the case

### tcps

As a tcp server to wait for client connections, reverse tcp mode of operation

Note: Only one client can be connected at a time.
multiple clients connected at the same time will wait until the previous close before handling new connections

### exec

exec system command. **NOTE: This is very dangerous!**

```yaml
  - type: exec
    params: "echo"
    i_rules:
      - regexp: '.*"method": ?"exec".*'
    o_rules:
      - regexp: '.*'
```

This `params` is command, example: `params: "/bin/sh"`

This Recvice JSONRPC, ignore `jsonrpc.method`. **jsonrpc.method is use of only match**

For example:

`params: "echo"`. jsonrpc2 `{"jsonrpc":"2.0","method":"exec","params":["-n", "xxx"],"id":"x"}`

Finally exec: `echo -n xxx`

### jsonrpc2

Is a simulator that can simulate jsonrpc correct and error messages, which is often used for development and debugging

The `params` field of this module is a bit more complicated: the

When it is json, this is inserted directly (only `result` and `error` are legal)

When it's not json, the result is assigned directly to `result`.

Simulation success: `params: 233` is the same as this.

```json
{"result": 233}
```

Simulation error:

```json
{"error": {"code": 0, "message": "xxxxx"}}
```

### logger

```yaml
params: "file:///tmp/ncp/test.log?size=128M&count=8&prefix=SB"
```

Why `//`, this is the standard way to write `//` when the url is part of `/tmp` is the root path

* size The maximum size of each log file, beyond which it will rotate
* count Keeps the last few files
* prefix log prefix

### mqtt

This module is a bit complicated so break it down to ncpio's mqtt and mqttd

Of course, you can also ignore this function and treat it as a normal io

All params are defined as strings, because yaml doesn't support save parsing, so it's not possible to do the same as `json.RawMessage`

The reason is: [go-yaml/issues#13](https://github.com/go-yaml/yaml/issues/13)

yaml is contextual, so it doesn't have the same context as json

So mqtt's params are specified as a path

## mqttd

![mqttd](docs/ncp-mqttd.svg)

Messages received through the IO module will arrive here

The config file [mqtt.yaml](conf/mqtt.yml)

### mqttd-rpc

config params:

```yaml
mqttd:
  rpc:
    qos: 0
    lru: 128
    i: "nodes/%s/rpc/recv"
    o: "nodes/%s/rpc/send"
```

There are two mqtt topics corresponding to `i` and `o`, sending and receiving. topic is full duplex, why do we need to split into two topics?

A simple use can be to set both values to be the same and use one topic for sending and receiving

* * *

Message reliability control, two modes. Let's review `tcp` and `kcp` first

`tcp` returns an acknowledgement message when it receives a message, and then sends a new message, which is acknowledged by a sliding window

Later, in the pursuit of low latency, `udp` is wrapped with reliability, and packets are sent violently to improve response time. This is what `kcp` does

The `kcp` approach has lower latency than `tcp`, but consumes more traffic than `tcp`.

Message reliability control, although `mqtt` but `qos` can be reliability control, set `qos 1` or `qos 2`

In the actual test product, `mqtt` itself `qos` does not perform well when the network conditions are bad, the real-time requirements are high, and the data volume is large

mqtt's `qos` will repeatedly confirm the arrival of packet messages (`qos` is doubled for broker, two clients), so use `qos 0` for reliability control at the upper layer itself

* * *

This introduces two parameters: `qos` and `lru`

Use mqtt's own `qos` (similar to `tcp` mode) to set `qos` to `2`. lru is set to `0`.

own reliability processing (`kcp`-like mode) to set `qos` to `0`. lru to a number other than `0` (default is `128`)

Why is there a `LRU` here?

So use a new mode, brute force data sending, using LRU to filter duplicate messages sent

It's more complicated to write your own message validation part to control reliability.

It is recommended that the demo use mqtt's own qos. product environment uses `LRU` and `qos 0` for this working mode

### mqttd-status

`status` & `network` is send self status

`status` send self status: mqtt last will message, network error

```json
{
  "msg":"neterror",
  "timestamp":"1619164695",
  "status":{
    "lat":"22.6876423001",
    "lng":"114.2248673001",
    "alt":"10088.0001"
  }
}
```

`status` config is `mqttd.static` . An extended field that can hold some meta information

`msg`: `online`, `offline`, `neterror`

### mqttd-network

```json
{"loss":0,"delay":5}
```

`network` send self and mqtt broker network status

* `loss` (1~100)%
* `delay` unit `ms`

### mqttd-trans

The received message will be judged as non-jsonrpc, if it is jsonrpc message will connect to mqtt broker in rpc way

If it is not jsonrpc, it will go through this module, which is a one-way data

tran socket recv

```json
{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}}
```

mqtt send topic: `nodes/:id/msg/weather`

```json
{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}
```

Automatically maps the outermost key of the json as part of the topic.
(Multiple outer keys will be split into multiple messages)

`gtran.prefix` config mqtt prefix

`trans` config mqtt topic `qos` and `retain`

For example:

```json
{
  "weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780},
  "battery":{"status":"ok"}
}
```

mqtt send topic: `nodes/:id/msg/weather`

```json
{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}
```

mqtt send topic: `nodes/:id/msg/battery`

```json
{"status":"ok"}
```

### mqttd-buildin

build-in `history` method

`history` get `trans` send history data

params:

Field | Type | Description
----- | ---- | -----------
topic | string | Topic
time  | string | such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h" [golang time#ParseDuration](https://golang.org/pkg/time/#ParseDuration)

Request:

```json
{
  "jsonrpc":"2.0",
  "id":"sdwc.1-1553321035000",
  "method":"history",
  "params":{
    "topic":"msg/weather",
    "time": "10s"
  }
}
```

Response:

```json
{
  "jsonrpc": "2.0",
  "id":"sdwc.1-1553321035000",
  "result":[
    {"1553321035": {"WD": 20, "WS": 5}},
    {"1553321034": {"WD": 10, "WS": 1}},
    {"1553321000": {"WD": 1, "WS": 1}}
  ]
}
```

The result[key] , key is timestamp

`ncp_online` and `ncp_offline` control online status

Is jsonrpc notification. not result, not params

## Example

If use: `tcpc`, `tcps`, `mqtt`, `logger`

Use this config.yaml :

```yaml
# type: tcps / tcpc / mqtt / logger / jsonrpc2 / api
ncpio:
  - type: tcpc
    params: "localhost:8980"
    i_rules:
      - regexp: '.*"method": ?"(webrtc|history|ncp)".*'
        invert: true
      - regexp: '.*"method".*'
    o_rules:
      - regexp: '.*'
  - type: tcps
    params: "localhost:1208"
    i_rules:
      - regexp: '.*"method": ?"webrtc".*'
    o_rules:
      - regexp: '.*"method".*'
        invert: true
      - regexp: '.*'
  - type: mqtt
    params: "config.yml"
    i_rules:
      - regexp: '.*'
    o_rules:
      - regexp: '.*'
  - type: logger
    params: "file:///tmp/ncp/test.log?size=128M&count=8&prefix=SB"
    i_rules:
      - regexp: '.*'
```

broker to mqttd send:

```json
{"jsonrpc":"2.0","id":"test.0-1553321035000","method":"webrtc","params":[]}
```

1. `mqtt` will first go through its own `o_rules` and find that it is `. *` (matching any message), so it broadcasts this message to `tcpc` `tcps` `logger`
2. `logger` is a logger: `i_rules` is usually `. *` writes to the log after receiving this message
3. `tcpc`'s `i_rules` matches `. *"method": ?" (webrtc|history|ncp)". *`, but `invert` reverses the match. This message is rejected
4. `i_rules` of `tcps` just matches this command. This command is sent to `tcps`
5. `tcps` receives the message filtered by `o_rules` and sends it to the message center, which broadcasts the message to every enabled `IO` other than itself (excluding the sender of the message).
6. `mqtt` also receives the message, and after receiving the broadcast message, it uses its own `i_rules` to determine if the message is ready to be sent.
7. `logger` will also receive this message, and after receiving the broadcast message, it will use its own `i_rules` to determine if the message needs to be logged to the log file

## TODO

* [ ] io udp's io maybe needs to be added
* [ ] io mqtt or need a configuration item to enable to use json and jsonrpc

