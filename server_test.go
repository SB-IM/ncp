package main

import (
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func fakeLogger() *log.Logger {
	_, logfile, _ := os.Pipe()
	return log.New(logfile, "[Server Test] ", log.LstdFlags)
}

func Test_GetMethods(t *testing.T) {
	server, client := net.Pipe()

	method_group_1 := []string{"test1", "test2"}
	method_group_2 := []string{"test3", "test4"}

	socketServer := &SocketServer{
		logger: fakeLogger(),
		links: []*Link{
			&Link{
				conn:    &server,
				methods: method_group_1,
			},
			&Link{
				conn:    &client,
				methods: method_group_2,
			},
		},
	}

	method_group_all := append(method_group_1, method_group_2...)

	methods := socketServer.GetMethods()

	if !StringSliceReflectEqual(methods, method_group_all) {
		t.Errorf("Has Method No Match")
		for _, method := range methods {
			t.Errorf(method)
		}
	}
}

func StringSliceReflectEqual(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}

func Test_getMethodMatchConns(t *testing.T) {
	server, client := net.Pipe()

	method_group_1 := []string{"test1", "test2"}
	method_group_2 := []string{"test3", "test4"}

	socketServer := &SocketServer{
		logger: fakeLogger(),
		links: []*Link{
			&Link{
				conn:    &server,
				methods: method_group_1,
			},
			&Link{
				conn:    &client,
				methods: method_group_2,
			},
		},
	}

	call_method := "test1"
	conns := socketServer.getMethodMatchConns(call_method)

	ch := make(chan string)
	go func() {
		buf := make([]byte, 4096)
		cnt, err := client.Read(buf)
		if err != nil || cnt == 0 {
			t.Errorf(err.Error())
		}
		ch <- strings.TrimSpace(string(buf[0:cnt]))
	}()

	for _, conn := range conns {
		(*conn).Write([]byte(call_method + "\n"))
	}

	if call_method != <-ch {
		t.Errorf("Method No Match")
	}
}

func Test_AddDelLink(t *testing.T) {
	server, client := net.Pipe()

	method_group_1 := []string{"test1", "test2"}
	method_group_2 := []string{"test3", "test4"}

	link_server := &Link{
		conn:    &server,
		methods: method_group_1,
	}
	link_client := &Link{
		conn:    &client,
		methods: method_group_2,
	}

	socketServer := &SocketServer{
		logger: fakeLogger(),
		links:  []*Link{},
	}

	socketServer.addLink(link_client)
	socketServer.addLink(link_server)

	link_client.methods = []string{"test5", "test6", "test7"}
	socketServer.addLink(link_client)

	socketServer.delLink(link_server)

	for _, link := range socketServer.links {
		if link != link_client {
			t.Errorf("Link not exist in []links")
		}
	}
}

func Test_SocketServerRecv(t *testing.T) {
	server, client := net.Pipe()
	socketServer := &SocketServer{
		logger: fakeLogger(),
	}

	output := make(chan string)
	go socketServer.recv(&server, output)

	reg_notify := `{"jsonrpc":"2.0","method":"reg","params":["test", "status", "webrtc"]}`
	test_msg := `{"test":"233"}`

	client.Write([]byte(reg_notify + "\n"))
	client.Write([]byte(test_msg + "\n"))

	if test_msg != <-output {
		t.Errorf("msg not match")

		for i, link := range socketServer.links {
			t.Errorf("-- conn is: " + strconv.Itoa(i))
			for _, method := range link.methods {
				t.Errorf(method)
			}
		}
	}
}

func Test_SocketServerRecvReg(t *testing.T) {
	server, client := net.Pipe()
	socketServer := &SocketServer{
		logger: fakeLogger(),
	}

	output := make(chan string)
	reg_request := `{"jsonrpc":"2.0","id":"test.0-00000000","method":"reg","params":["test", "status", "webrtc"]}`
	reg_result := `{"jsonrpc":"2.0","result":["test","status","webrtc"],"id":"test.0-00000000"}`

	no_output := make(chan string)
	go socketServer.recv(&server, no_output)

	go func() {
		buf := make([]byte, 4096)
		cnt, err := client.Read(buf)
		if err != nil || cnt == 0 {
			t.Errorf("IO error")
		}
		output <- strings.TrimSpace(string(buf[0:cnt]))
	}()

	client.Write([]byte(reg_request + "\n"))
	result := <-output

	if result != reg_result {
		t.Errorf("Result Not Match")
		t.Errorf(result)
	}
}

func Test_SocketServerSend(t *testing.T) {
	server, client := net.Pipe()
	server2, client2 := net.Pipe()

	method_group_1 := []string{"test", "test2"}
	method_group_2 := []string{"test3", "test4"}

	socketServer := &SocketServer{
		logger: fakeLogger(),
		links: []*Link{
			&Link{
				conn:    &server,
				methods: method_group_1,
			},
			&Link{
				conn:    &server2,
				methods: method_group_2,
			},
		},
	}

	/*
		socketServer := &SocketServer{}

		// method registered
		reg_request := `{"jsonrpc":"2.0","id":"test.0-00000000","method":"reg","params":["test", "status", "webrtc"]}`
		reg_request2 := `{"jsonrpc":"2.0","id":"test.0-00000000","method":"reg","params":["test2"]}`

		output := make(chan string)
		go socketServer.recv(&server, output)
		go socketServer.recv(&server2, output)
		client.Write([]byte(reg_request + "\n"))
		client2.Write([]byte(reg_request2 + "\n"))

		// Need Wait reg_request2 sueccess
		time.Sleep(100 * time.Millisecond)
		// registered done
	*/

	input := make(chan string)
	jsonrpc_request := `{"jsonrpc":"2.0","id":"test.0-00000000","method":"test","params":["status", "webrtc"]}`

	go socketServer.send(input)
	input <- jsonrpc_request

	buf := make([]byte, 4096)
	cnt, err := client.Read(buf)
	if err != nil || cnt == 0 {
		t.Errorf("IO error")
	}
	msg := strings.TrimSpace(string(buf[0:cnt]))
	if jsonrpc_request != msg {
		t.Errorf("recv success")
	}

	// Should not receive data

	// set SetReadDeadline
	err = client2.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	if err != nil {
		t.Errorf("SetReadDeadline failed: " + err.Error())
	}

	cnt, err = client2.Read(buf)
	if err == nil || cnt != 0 {
		t.Errorf("IO error")
	}

}
