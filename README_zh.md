# Node control protocol

[README](README.md) | [中文文档](README_zh.md)

## What is ncp

这是通信模块，并没有规定具体用途，可以用来做很多事

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Compile && Run](#compile--run)
* [Architecture design](#architecture-design)
* [NcpIOs](#ncpios)
  * [Rules](#rules)
  * [Debug](#debug)
  * [Name](#name)
* [IO 模块](#io-模块)
  * [API](#api)
  * [tcpc](#tcpc)
  * [tcps](#tcps)
  * [jsonrpc2](#jsonrpc2)
  * [logger](#logger)
  * [mqtt](#mqtt)
* [mqttd](#mqttd)
  * [mqtt-rpc](#mqtt-rpc)
  * [mqtt-status](#mqtt-status)
  * [mqtt-tran](#mqtt-tran)
  * [Q&A](#qa)
* [Example](#example)
* [TODO](#todo)

<!-- vim-markdown-toc -->

## Compile && Run

go >= 1.13

```sh
make

# run
cp conf/config-dist.yml config.yml
ncp -c config.yml
```

## Architecture design

```sh
Cloud -------------- Mqtt broker
                      |
                      |
                      |
jsonrpc2 simulation -- ncp --------- logger
                    /      \
                   /        \
                  /          \
           socket client      socket server
           jsonrpc            jsonrpc
          airctl/dockctl      startlight

```

## NcpIOs

这个是整个应用的核心，也可以当库来使用

原本定义了七大 IO 模块。还有个叫 `webui` 的模块。

这是是调试使用的，ncp 原本是给 iot 设备设计的网关程序

现场部署调试设备时经常无法连接云端来调试，所以还有 `webui` 的调试专用模块还没有进行开发

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

所以有三条日志（RECV, BROADCAST, SEND）顺序时这样的：

1. RECV: `IO` 刚收到消息时
2. 经过 `o_rules`
3. BROADCAST: 向全体广播消息时
4. 经过 `i_rules`
5. SEND: `IO` 可以发送消息时

### Name

这每种类型当模块时可以存在多个的。可以给每个都命个名字

默认情况下会自动生成一个唯一的，从 `0`, `1`, `2` 依次递增

建议不要设置，如果设置了相同 `name` 字段，在广播消息时是通过 `name` 来放在自己发送给自己的。

相同 `name` 会收不到同名广播信息。**请确认好，是有意不接受广播，还是不小心设置成了一样的**

## IO 模块

`IO` 模块有这六种类型：

Type                 | Description
-------------------- | -----------
[API](#api)          | build-in interface
[tcpc](#tcpc)        | TCP socket client
[tcps](#tcps)        | TCP socket server
[jsonrpc2](#jsonrpc2)| Connect Mqtt Broker
[logger](#logger)    | TCP socket server
[mqtt](#mqtt)        | Connect Mqtt Broker

核心是六个 `IO` 模块，通过配置文件来定义不同的组合来实现不同的功能

为什么没有 `udp server` 和 `udp client` ?

原本开发的时候并没有这个需求，有些消息需要可达性验证。 `udp` 并不合适

### API

是一个特殊接口

当把 ncpio 作为库使用时才会用到，通过这个 API 是个 IO 模块，可以直接通信

注意：单元测试就是使用这个接口测试的

### tcpc

作为一个 tcp client 去连接 server，最普通的工作模式

tcpc, tcps 消息分割使用 `\n` 作为分隔符，会自动 `chomp` `\r\n` 当情况

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

* size 每个日志人家最大体积，超过这个体积会 rotate
* count 保留最近几个文件
* prefix 日志前缀

### mqtt

mqtt 和 json, jsonrpc2.0 绑定了，当然也可以忽略这个功能当成普通的 io 来处理

全部的 params 定义成字符串类型，因为 yaml 不支持保存解析，没法像 `json.RawMessage` 一样

原因是：[go-yaml/issues#13](https://github.com/go-yaml/yaml/issues/13)

yaml 是有上下文的，所以没有像 json 一样

因此 mqtt 的 params 是指定的一个路径

## mqttd

这个模块比较复杂所以分解来说

### mqtt-rpc

有两个对应 topic , 发送和接收

为什么这里会有 `LRU` ?

消息可达性控制，虽然 `mqtt` 但 `qos` 可以进行可达性控制，设置 `qos 1` 或 `qos 2`

实测产品情况，网络条件恶劣，实时要求高，数据量大时，`mqtt` 本身 `qos` 表现并不理想

所以使用了一种全新的模式，暴力发送数据，使用 LRU 来过滤重复发送的消息

这两个模式有点像 `tcp` 和 `kcp`。本程序作者建议使用 `LRU` 和 `qos 0` 的这种工作模式

集成了 `history` 命令，和 `jsonrpc` 绑定在一起。当然可以忽略这个消息

### mqtt-status

`status` 和 `network` 是特殊定义的两个 topic 模块。

`status` 会发送一下本身的状态，比如遗嘱消息

`network` 会上报自己的网络状态

### mqtt-tran

type: `mqtt` build-in

tran socket recv

```json
{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}}
```

mqtt send topic: `node/:id/msg/weather`

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

mqtt send topic: `node/:id/msg/weather`

```json
{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}
```

mqtt send topic: `node/:id/msg/battery`

```json
{"status":"ok"}
```

### Q&A

jsonrpc 的 topic 是全双工的，为什么还要分成两个 topic ？

简单的使用完全可以使用把这两个值设置成相同的

但配合 retain 会有一下特殊的的需求和用法，可以自行去探索

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

