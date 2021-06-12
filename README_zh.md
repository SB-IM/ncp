# Node control protocol

Edge gateway core communication service

[![Build Status](https://github.com/SB-IM/ncp/workflows/ci/badge.svg)](https://github.com/SB-IM/ncp/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/SB-IM/ncp)](https://goreportcard.com/report/github.com/SB-IM/ncp)
[![Go Reference](https://pkg.go.dev/badge/github.com/SB-IM/ncp.svg)](https://pkg.go.dev/github.com/SB-IM/ncp)
[![GitHub release](https://img.shields.io/github/tag/SB-IM/ncp.svg?label=release)](https://github.com/SB-IM/ncp/releases)
[![license](https://img.shields.io/github/license/SB-IM/ncp.svg?maxAge=2592000)](https://github.com/SB-IM/ncp/blob/master/LICENSE)

[README](README.md) | [中文文档](README_zh.md)

## What is ncp

这是通信模块，并没有规定具体用途，可以用来做很多事

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

这个是整个应用的核心，也可以当库来使用

原本定义了七大 IO 模块。还有个叫 `webui` 的模块。

这是是调试使用的，ncp 原本是给 iot 设备设计的网关程序

现场部署调试设备时经常无法连接云端来调试，所以还有 `webui` 的调试专用模块，不过这个还没有进行开发

这里使用的是一种 `IO` 模块和消息中心的模式，每一个 `IO` 模块收到的消息和都会广播给其他的 `IO` 模块

通过 `rules` 来判断这条消息是否可以接收，当没有匹配任何规则时为否定（不接受这条消息）

这个模块可以作为库使用。消息可直接使用 `api` 接口来发送和接收

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

有两种：

* `o_rules` output the io rule (after recv)
* `i_rules` input the io rule (before send)

1. `i_rules` 和 `o_rules` 是通过正则匹配的，所以传输内容要求明文
2. 规则可以有多条，依此匹配，只要符合其中的一条，就可以通过规则
3. `invert` 字段可以控制正则匹配结果是非要取反
4. 当没有匹配任何规则时，消息不会通过

### Debug

rules 使用正则表达式部分的复杂，很难调试。所以增加了调试模式

```bash
./ncp -debug
```

每条信息通过 `IO` 模块至少要经过两个 rules 组

所以有三条日志（RECV, BROADCAST, SEND）顺序是这样的：

1. `RECV`: `IO` 刚收到消息时
2. 通过 `o_rules`
3. `BROADCAST`: 向全体广播消息时
4. 通过 `i_rules`
5. `SEND`: `IO` 可以发送消息时

### Name

每种类型的 `IO` 模块时可以存在多个的。可以给每个都命个名字

默认情况下会自动生成一个唯一的名字，从 `0`, `1`, `2` 依次递增

建议手动不要设置 `name` 字段，如果设置了相同 `name` 字段，在广播消息时是通过 `name` 来放在自己发送给自己的。

相同 `name` 会收不到同名广播信息。**请确认好，是有意不接受广播，还是不小心设置成了一样的**

## IOs

`IO` 模块有这六种类型：

Type                 | Description
-------------------- | -----------
[API](#api)          | build-in interface
[tcpc](#tcpc)        | TCP socket client
[tcps](#tcps)        | TCP socket server
[jsonrpc2](#jsonrpc2)| JSONRPC 2.0 simulation
[logger](#logger)    | Record message log
[mqtt](#mqtt)        | Connect Mqtt Broker

核心是六个 `IO` 模块，通过配置文件来定义不同的组合来实现不同的功能

为什么没有 `udp server` 和 `udp client` ?

原本开发的时候并没有这个需求，有些消息需要可靠性验证。 `udp` 并不合适

### API

是一个特殊接口

当把 ncpio 作为库使用时才会用到，通过这个 API 是个 IO 模块，可以直接通信

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

注意：单元测试就是使用这个接口测试的

API 接口最多只能有一个

### tcpc

作为一个 tcp client 去连接 server，最普通的工作模式

tcpc, tcps 消息分割使用 `\n` 作为分隔符，会自动 `chomp` `\r\n` 的情况

### tcps

作为一个 tcp server 去等待 client 连接，反向 tcp 工作模式

注意：同一时刻只能连接一个 client。多个 client 同时连接会等到之前的 close 后再处理新连接

### jsonrpc2

是一个模拟器，可以模拟 jsonrpc 正确和错误消息，开发和调试时会经常使用这个

这个模块的 `params` 字段稍微有点复杂：

params 会尝试是否为 json。为 json 时，会把把这个直接插进来（只有 `result` 和 `error` 时合法）

当不为 json 是会把这个结果直接赋给 `result`

模拟成功：`params: 233` 和这个等同：

```json
{"result": 233}
```

模拟失败：

```json
{"error": {"code": 0, "message": "xxxxx"}}
```

### logger

params 为 url 格式

```yaml
params: "file:///tmp/ncp/test.log?size=128M&count=8&prefix=SB"
```

为啥 `///` ，这个是标准写法 `//` 时 url 的一部分 `/tmp` 是根路径

* size 每个日志文件最大体积，超过这个体积会 rotate
* count 保留最近几个文件
* prefix 日志前缀

### mqtt

这个模块比较复杂所以分解来说 ncpio 的 mqtt 和 mqttd

当然也可以忽略这个功能当成普通的 io 来处理

全部的 params 定义成字符串类型，因为 yaml 不支持保存解析，没法像 `json.RawMessage` 一样

原因是：[go-yaml/issues#13](https://github.com/go-yaml/yaml/issues/13)

yaml 是有上下文的，所以没有像 json 一样

因此 mqtt 的 params 是指定的一个路径

## mqttd

![mqttd](docs/ncp-mqttd.svg)

通过 IO 模块收到的消息会到达这里

这个模块和 json, jsonrpc2.0 绑定了，和业务需有是有一定耦合的

这个设置文件 [mqtt.yaml](conf/mqtt.yml)

### mqttd-rpc

可以配置参数：

```yaml
mqttd:
  rpc:
    qos: 0
    lru: 128
    i: "nodes/%s/rpc/recv"
    o: "nodes/%s/rpc/send"
```

有两个对应 `i` 和 `o` mqtt topic , 发送和接收。topic 是全双工的，为什么还要分成两个 topic ？

简单的使用完全可以使用把这两个值设置成相同的，收发使用一个 topic

* * *

消息可靠性控制，两种模式。先来回顾一下 `tcp` 和 `kcp`

`tcp` 收到消息返回确认消息，再发送新消息，通过滑动窗口来确认

后来为了追求低延时，把 `udp` 加上可靠性封装，暴力发送数据包来提高响应速度。就是 `kcp` 的做法

`kcp` 的做法延时比 `tcp` 低，但会比 `tcp` 多消耗更多的流量

消息可靠性控制，虽然 `mqtt` 但 `qos` 可以进行可靠性控制，设置 `qos 1` 或 `qos 2`

实测产品情况，网络条件恶劣，实时要求高，数据量大时，`mqtt` 本身 `qos` 表现并不理想

mqtt 的 `qos` 会反复确认包消息到达（`qos` 是对于 broker ，两个客户端会翻倍），所以使用 `qos 0` 在上层自己进行可靠性控制

* * *

这里面引入两个参数：`qos` 和 `lru`

使用 mqtt 自己的 `qos` （类似 `tcp` 的模式）要把 `qos` 设置成 `2`。lru 设置成 `0`

自己的可靠性处理（类似 `kcp` 的模式）要把 `qos` 设置成 `0`。lru 设置成 `0` 以外的数（默认是 `128`）

为什么这里会有 `LRU` ?

所以使用了一种全新的模式，暴力发送数据，使用 LRU 来过滤重复发送的消息

自己控制可靠性要自行写消息验证部分，比较复杂。

建议 demo 使用 mqtt 自己的 qos。产品环境使用 `LRU` 和 `qos 0` 的这种工作模式

### mqttd-status

`status` 和 `network` 是特殊定义的两个 topic 模块。

`status` 会发送一下本身的状态，比如遗嘱消息。由于网络异常发送遗嘱消息

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

`status` 对应的是配置 `mqttd.static` 里面的东西。相当于一个扩展字段，可以放了一些元信息

用于区分：正常连接，正常关闭，网络异常的状态

### mqttd-network

```json
{"loss":0,"delay":5}
```

`network` 会上报自己和 mqtt broker 之间的网络状态

* `loss` 是丢包率
* `delay` 是延时（单位 `ms`）

### mqttd-trans

这是一个内置的模块，收到的消息会判断是非为 jsonrpc ，如果是 jsonrpc 消息会以 rpc 方式连接 mqtt broker

如果不是 jsonrpc 就会经过这个模块，这里面数据是单向的

tran socket recv

```json
{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}}
```

mqtt send topic: `nodes/:id/msg/weather`

```json
{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}
```

会自动化把 json 的最外层的键映射为 topic 的一部分。（多个外层键会分割成多条消息）

`gtran.prefix` 可以配置使用什么样的前缀

`trans` 可以对于每一个 topic 预先定义参数 `qos` 或 `retain`

比如：

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

这个模块内置了几条命令

集成了 `history` 命令，和 `jsonrpc` 绑定在一起。当然可以忽略这个消息

`history` 可以返回 `trans` 收到的历史消息

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

`ncp_online` 和 `ncp_offline` 可以控制在线状态（声明这个设备关闭并不等于通信全部关闭）

这两条命令是 jsonrpc notification。没有返回，也没有参数

## Example

比如使用 `tcpc`, `tcps`, `mqtt`, `logger`

使用这样一个配置文件：

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

通过 broker 给 mqtt 模块发送个

```json
{"jsonrpc":"2.0","id":"test.0-1553321035000","method":"webrtc","params":[]}
```

1. `mqtt` 会先通过自己的 `o_rules`，发现是 `.*` （匹配任何消息），因此向 `tcpc` `tcps` `logger` 广播这条消息
2. `logger` 是个日志记录器：`i_rules` 一般都是 `.*` 收到这个消息后写入日志
3. `tcpc` 的 `i_rules` 第一条就匹配到了 `.*"method": ?"(webrtc|history|ncp)".*` ，但 `invert` 反转匹配结果。这条消息拒绝接收
4. `tcps` 的 `i_rules` 刚好匹配这条命令。这个命令发送给 `tcps`
5. `tcps` 收到消息经过 `o_rules` 过滤后发送给消息中心，消息中心会像自己以外的每一个启用的 `IO` 广播这条消息（不包括消息的发送者）。
6. `mqtt` 也会收到这条消息，收到广播消息之后通过自己的 `i_rules` 来判断这条消息是否可以发送出去
7. `logger` 也会收到这条消息，收到广播消息之后通过自己的 `i_rules` 来判断这条消息是非需要记录到日志文件里

## TODO

* [ ] io udp 的 io 看情况决定是非需要加入
* [ ] io mqtt 或者需要个配置项来决定使用要耦合 json 和 jsonrpc

