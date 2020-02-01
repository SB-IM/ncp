package main

import (
	"encoding/json"

	"github.com/SB-IM/jsonrpc2"
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
	rpc := jsonrpc2.Jsonrpc{}
	//err := json.Unmarshal(req, &jsonrpc)
	json.Unmarshal([]byte(s), &rpc)
	if rpc.IsNotify() {
		return true
	}


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

func linkCall(req []byte, id string) ([]byte, error, func([]byte) ([]byte, error)) {
	jsonrpc_req := jsonrpc2.WireRequest{}
	err := json.Unmarshal(req, &jsonrpc_req)
	src_id := *jsonrpc_req.ID

	callback := func(res []byte) ([]byte, error) {
		jsonrpc_res := jsonrpc2.WireResponse{}
		err := json.Unmarshal([]byte(res), &jsonrpc_res)
		if err != nil {
			return []byte(""), err
		}
		jsonrpc_res.ID = &src_id
		return json.Marshal(jsonrpc_res)
	}

	if err != nil {
		return []byte(""), err, callback
	}

	raw_params, _ := jsonrpc_req.Params.MarshalJSON()
	var params []string
	err = json.Unmarshal(raw_params, &params)
	if err != nil {
		return []byte(""), err, callback
	}

	jsonrpc_req.Method = params[0]
	jsonrpc_req.Params = nil
	jsonrpc_req.ID.Name = id

	jsonrpc, err := json.Marshal(jsonrpc_req)
	return jsonrpc, err, callback
}

func detachTran(raw []byte) (map[string][]byte) {
	srcMap := make(map[string]*json.RawMessage)
	dstMap := make(map[string][]byte)
	json.Unmarshal(raw, &srcMap)
	for k, v := range srcMap {
		dstMap[k], _ = v.MarshalJSON()
	}
	return dstMap
}

func getReg(raw []byte) ([]string, bool) {
	var params []string
	if getJSONRPC(string(raw)).Method == "reg" {
		raw_params, _ := getJSONRPC(string(raw)).Params.MarshalJSON()
		_ = json.Unmarshal(raw_params, &params)
		return params, true
	}
	return params, false
}

