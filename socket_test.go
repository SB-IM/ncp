package main

import (
	"net"
	"strings"
	"testing"
)

func Test_SocketClient_record(t *testing.T) {
	socketClient := &SocketClient{
		logger: fakeLogger(),
	}
	socketClient.record([]byte(test_jsonrpc_send))
	if socketClient.running == nil {
		t.Errorf("Should has record")
	}

	socketClient.record([]byte(`{"jsonrpc":"2.0","result":["status","webrtc"],"id":"test.0-00000000"}`))
	if socketClient.running == nil {
		t.Errorf("Should has record")
	}

	socketClient.record([]byte(test_jsonrpc_recv_r))
	if socketClient.running != nil {
		t.Errorf(string(*socketClient.running))
	}
}

func Test_SocketClient_recv(t *testing.T) {
	server, client := net.Pipe()
	socketClient := &SocketClient{
		logger: fakeLogger(),
	}

	output := make(chan string)
	go socketClient.recv(server, output)

	test_msg := `{"test":"233"}`

	client.Write([]byte(test_msg + "\n"))
	if msg := <-output; test_msg != msg {
		t.Errorf("msg not match")
		t.Errorf(msg)
	}

	client.Write([]byte(test_msg + "\n" + test_msg + "\n"))
	if msg := <-output; test_msg != msg {
		t.Errorf("msg not match")
		for _, v := range strings.Split(msg, "\n") {
			t.Errorf(v)
		}
	}
}
