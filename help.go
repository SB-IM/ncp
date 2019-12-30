package main

import (
  "encoding/json"
	"strconv"
	"time"
)

func isJSON(s string) bool {
  var js map[string]interface{}
  return json.Unmarshal([]byte(s), &js) == nil
}

type JSONRPC struct {
  Jsonrpc string
  Method string
	Params     *json.RawMessage `json:"params,omitempty"`
  Result interface{}
  Error interface{}
  Id string
}

func getJSONRPC(s string) JSONRPC {
  jsonrpc := JSONRPC{}
  json.Unmarshal([]byte(s), &jsonrpc)
  return jsonrpc
}

func isJSONRPCSend(s string) bool {
  return getJSONRPC(s).Method != ""
}

func isJSONRPCRecv(s string) bool {
  return getJSONRPC(s).Result != nil || getJSONRPC(s).Error != nil
}

type RpcRun struct {
	run []string
}

func (this *RpcRun) Run(s string) bool {
	id := getJSONRPC(s).Id
	if func(str string, array []string) bool {
		for _, r := range array {
			if r == id {
				return true
			}
		}
		return false
	}(id, this.run) {
		return false
	} else {
		if len(this.run) >= 128 {
			this.run = this.run[1:]
		}
		this.run = append(this.run, id)
		return true
	}
}

func confirmNotice(s string) string {
	return `{"jsonrpc": "2.0", "method": "ack", "params": { "id": "` + getJSONRPC(s).Id + `" }}`
}

func isNcp(s string) bool {
	method := getJSONRPC(s).Method
	isncp := false

	for _, m:= range []string{"ncp", "status", "upload", "download", "shell", "webrtc"} {
		if method == m { isncp = true }
	}
  return isJSONRPCSend(s) && isncp
}

func isLink(s string) bool {
	if getJSONRPC(s).Method == "link" {
		return true
	} else {
		return false
	}
}

func linkCall(raw string, caller int) string {
	rpc := getJSONRPC(raw)
	var params []string
	json.Unmarshal(*rpc.Params, &params)

	bit13_timestamp := string([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))[:13])
	return `{"jsonrpc":"2.0","id":"ncp.` + strconv.Itoa(caller) + `-` + bit13_timestamp + `","method":"` + params[0] + `","params":[]}`
}

