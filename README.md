# Node control protocol

#### Compile && Run
go >= 1.12.4

```sh
make

# run
cp conf/config-dist.yml config.yml
ncp -c config.yml
```

### Architecture design

```
Cloud -------------- Mqtt broker
  |                   |
  | httpUpload        |
  |                   |
built-in commond  --- ncp --------- tran json
                    /      \
                   /        \
                  /          \
           socket client      socket server
           jsonrpc            jsonrpc
          airctl/dockctl      has 'reg' commond

```

Name | Description
---- | -----------
built-in | eg: `status, upload, download`
socket client | only: jsonrpc 2.0, legacy, plan remove
socket server | only: jsonrpc 2.0, reg method wait call
tran | Single status message transmission

#### reg
```json
{"jsonrpc": "2.0", "method": "reg", "params": ["webrtc", "webrtc_close"], "id": 1}
```

#### tran

tran socket recv
```json
{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}}
```
mqtt send topic: `node/:id/msg/weather`

```json
{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}
```

