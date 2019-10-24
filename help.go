package main

import (
  "encoding/json"
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

func isNcp(s string) bool {
	method := getJSONRPC(s).Method
	isncp := false

	for _, m:= range []string{"ncp", "status", "upload", "download", "shell", "webrtc"} {
		if method == m { isncp = true }
	}
  return isJSONRPCSend(s) && isncp
}
