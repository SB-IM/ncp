# Node control protocol

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

Type | Description
---- | -----------
api  | build-in interface
tcpc | TCP socket client
tcps | TCP socket server
mqtt | Connect Mqtt Broker
logger  | Record message log
jsonrpc2| JSONRPC 2.0 simulation

### Rule

`i_rules` input the io rule (before send)

`o_rules` output the io rule (after recv)

- default matched false: disallow
- if regexp matched: allow
- `invert` Invert matched result: !allow
- pipeline matched

### mqtt tran

type: `mqtt` build-in

tran socket recv

```json
{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}}
```

mqtt send topic: `node/:id/msg/weather`

```json
{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}
```
