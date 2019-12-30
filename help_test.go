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
	//t.Errorf(linkCall(linkcall, 2))

	if getJSONRPC(linkCall(linkcall, 2)).Method != "power_on_drone" {
		t.Errorf("Not Link")
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
