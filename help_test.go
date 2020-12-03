package main

import (
	"encoding/json"
	"strconv"
	"testing"
)

var test_jsonrpc_send = `{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`
var test_jsonrpc_recv_r = `{"jsonrpc":"2.0","result":["move_door_open","power_ir_on"],"id":"sdwc.1-1553321035000"}`
var test_jsonrpc_recv_e = `{"jsonrpc":"2.0","error":["move_door_open","power_ir_on"],"id":"sdwc.1-1553321035000"}`
var test_jsonrpc_send_ncp = `{"jsonrpc":"2.0","id":"sdwc.1-155332103904","method":"ncp","params":["status"]}`

func Test_isJSON(t *testing.T) {

	if isJSON("aaaaaa") {
		t.Errorf("Not JSON")
	}

	if !isJSON(test_jsonrpc_send) {
		t.Errorf("Is JSON")
	}
}

func Test_isJSONRPCSend(t *testing.T) {

	if !isJSONRPCSend(test_jsonrpc_send) {
		t.Errorf("Not JSON Send")
	}

	if isJSONRPCSend(test_jsonrpc_recv_r) {
		t.Errorf("Not JSON Send")
	}

	if isJSONRPCSend(test_jsonrpc_recv_e) {
		t.Errorf("Not JSON Send")
	}

	if isJSONRPCSend("aaaaa") {
		t.Errorf("Not JSON")
	}
}

func Test_isJSONRPCRecv(t *testing.T) {

	if !isJSONRPCRecv(test_jsonrpc_recv_r) {
		t.Errorf("Not JSON Recv")
	}

	if !isJSONRPCRecv(test_jsonrpc_recv_e) {
		t.Errorf("Not JSON Recv")
	}

	if isJSONRPCRecv(test_jsonrpc_send) {
		t.Errorf("Not JSON Recv")
	}

	if isJSONRPCRecv("aaaaa") {
		t.Errorf("Not JSON")
	}
}

func Test_RpcRun(t *testing.T) {
	m := RpcRun{}
	m.Run(test_jsonrpc_send)

	if m.Run(test_jsonrpc_send) {
		t.Errorf("No Duplicate Filter")
	}

	if !m.Run(test_jsonrpc_send_ncp) {
		t.Errorf("Excessive Filter")
	}

	if m.Run(test_jsonrpc_send_ncp) {
		t.Errorf("No Duplicate Filter")
	}
}

func Test_RpcRun_Notify(t *testing.T) {
	test_jsonrpc_notify := `{"jsonrpc":"2.0","method":"dooropen","params":[]}`
	m := RpcRun{}

	if !m.Run(test_jsonrpc_notify) {
		t.Errorf("Filter Notify")
	}

	if !m.Run(test_jsonrpc_notify) {
		t.Errorf("Filter Notify")
	}
}

func Test_RpcRun_limit(t *testing.T) {
	m := RpcRun{}
	m.Run(test_jsonrpc_send)

	if m.Run(test_jsonrpc_send) {
		t.Errorf("No Duplicate Filter")
	}

	// Max record 128
	for i := 0; i <= 127; i++ {
		call := `{"jsonrpc":"2.0","id":"` + strconv.Itoa(i) + `","method":"link","params":["power_on_drone"]}`
		m.Run(call)
	}

	if !m.Run(test_jsonrpc_send) {
		t.Errorf("Filter Max not 128 items")
	}
}

func Test_isNcp(t *testing.T) {

	if !isNcp(test_jsonrpc_send_ncp) {
		t.Errorf("Not Ncp")
	}

	if isNcp(test_jsonrpc_send) {
		t.Errorf("Not Ncp")
	}
}

func Test_isLink(t *testing.T) {
	linkcall := `{"jsonrpc":"2.0","id":"sdwc.1-155332103904","method":"link","params":["power_on_drone"]}`

	if !isLink(linkcall) {
		t.Errorf("Not Link")
	}

	if isLink(test_jsonrpc_send) {
		t.Errorf("Is Link")
	}
}

func Test_linkCall(t *testing.T) {
	linkcall := `{"jsonrpc":"2.0","id":"sdwc.1-155332103904","method":"link","params":["power_on_drone"]}`

	req, _, _ := linkCall([]byte(linkcall), "test.0-0000")
	if getJSONRPC(string(req)).Method != "power_on_drone" {
		t.Errorf("Not Link")
	}
}

func Test_linkCall_id(t *testing.T) {
	linkcall := `{"jsonrpc":"2.0","id":"test","method":"link","params":["power_on_drone"]}`

	test_jsonrpc_recv_r := `{"jsonrpc":"2.0","result":["move_door_open","power_ir_on"],"id":"sdwc.1-1553321035000"}`
	_, _, callback := linkCall([]byte(linkcall), "test.0-0000")
	res, _ := callback([]byte(test_jsonrpc_recv_r))
	if getJSONRPC(linkcall).Id != getJSONRPC(string(res)).Id {
		t.Errorf(getJSONRPC(linkcall).Id)
		t.Errorf(getJSONRPC(string(res)).Id)
		t.Errorf("id not equal")
	}
}

func Test_confirmNotice(t *testing.T) {
	id := "test.0-155332103904"
	linkcall := `{"jsonrpc":"2.0","id":"` + id + `","method":"link","params":["power_on_drone"]}`
	if rpc := getJSONRPC(confirmNotice(linkcall)); rpc.Method != "ack" {
		t.Errorf("Not Method: 'ack'")
	} else {
		var params struct {
			Id string `json:"id"`
		}
		json.Unmarshal(*rpc.Params, &params)

		if params.Id != id {
			t.Errorf(params.Id)
		}
	}
}

func Test_getReg(t *testing.T) {
	reg := `{"jsonrpc":"2.0","id":"test.0-00000000","method":"reg","params":["status", "webrtc"]}`

	methods, str := getReg([]byte(reg))

	if len(methods) == 0 {
		t.Errorf("Not Reg")
		for _, method := range methods {
			t.Errorf(method)
		}
	}

	result := `{"jsonrpc":"2.0","result":["status","webrtc"],"id":"test.0-00000000"}`
	if str != result {
		t.Errorf("Not Request")
	}
}

func Test_getRegNotify(t *testing.T) {
	reg := `{"jsonrpc":"2.0","method":"reg","params":["status", "webrtc"]}`

	methods, str := getReg([]byte(reg))

	if len(methods) == 0 {
		t.Errorf("Not Reg")
		for _, method := range methods {
			t.Errorf(method)
		}
	}

	result := ""
	if str != result {
		t.Errorf("Not Request")
	}
}

