package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func socketClient(addr string, logger *log.Logger, input chan string, output chan string) {
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(3 * time.Second)
		} else {
			logger.Println("New connect", &conn)
			go socketSend(conn, logger, input)
			socketRecv(conn, logger, output)
			logger.Println("Connect err try reconnect")
		}
	}
}

func socketRecv(conn net.Conn, logger *log.Logger, ch chan string) {
	buf := make([]byte, 4096)
	for {
		cnt, err := conn.Read(buf)
		if err != nil || cnt == 0 {
			logger.Println("Socket close")
			conn.Close()
			break
		}
		msg := strings.TrimSpace(string(buf[0:cnt]))
		logger.Println("Recv:", msg)
		ch <- msg
	}
}

func socketSend(conn net.Conn, logger *log.Logger, ch chan string) {
	for msg := range ch {
		logger.Println("Send:", msg)
		_, err := conn.Write([]byte(msg + "\n"))

		if err != nil {
			logger.Println("Error:", err)
			ch <- msg
			break
		}
	}
}
