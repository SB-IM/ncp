package main

import (
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

func Test_isNcp(t *testing.T) {

  if !isNcp(test_jsonrpc_send_ncp) {
    t.Errorf("Not Ncp")
  }

  if isNcp(test_jsonrpc_send) {
    t.Errorf("Not Ncp")
  }
}

