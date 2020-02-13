package main

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"

	"github.com/SB-IM/jsonrpc2"
)

type SocketClient struct {
	running *[]byte
	logger  *log.Logger
}

func (this *SocketClient) record(raw []byte) {
	rpc := jsonrpc2.Jsonrpc{}
	err := json.Unmarshal(raw, &rpc)
	if err != nil || rpc.IsNotify() {
		return
	}

	if !rpc.IsResponse() {
		this.running = &raw
		return
	}

	if this.running != nil {
		run_rpc := jsonrpc2.Jsonrpc{}
		json.Unmarshal(*(this.running), &run_rpc)
		if rpc.ID.String() == run_rpc.ID.String() {
			this.running = nil
		}
	}
}

func (this *SocketClient) Run(addr string, input chan string, output chan string) {
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(3 * time.Second)
		} else {
			this.logger.Println("New connect", &conn)

			if this.running != nil {
				this.logger.Println("ReSend:", string(*(this.running)))
				output <- string(*(this.running))
			}
			go this.send(conn, input)
			this.recv(conn, output)
			this.logger.Println("Connect err try reconnect")
		}
	}
}

func (this *SocketClient) recv(conn net.Conn, ch chan string) {
	buf := make([]byte, 4096)
	for {
		cnt, err := conn.Read(buf)
		if err != nil || cnt == 0 {
			this.logger.Println("Socket close")
			conn.Close()
			break
		}
		msg := strings.TrimSpace(string(buf[0:cnt]))
		for _, v := range strings.Split(msg, "\n") {
			this.logger.Println("Recv:", v)
			this.record([]byte(v))
			ch <- v
		}
	}
}

func (this *SocketClient) send(conn net.Conn, ch chan string) {
	for msg := range ch {
		this.logger.Println("Send:", msg)
		this.record([]byte(msg))
		_, err := conn.Write([]byte(msg + "\n"))

		if err != nil {
			this.logger.Println("Error:", err)
			ch <- msg
			break
		}
	}
}
