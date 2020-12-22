package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"sb.im/ncp/ncpio"
)

func socket() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config, err := getConfig("e2e/socket.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ncpios := ncpio.NewNcpIOs(config.NcpIO)
	go ncpios.Run(ctx)

	// wait tcp server startup
	time.Sleep(3 * time.Millisecond)

	msg1 := "2333333333333"
	msg2 := "4555555555555"

	ncpio.I <- []byte(msg1)
	ncpio.I <- []byte(msg2)

	for i := 0; i < 1000000; i++ {
		ncpio.I <- []byte(fmt.Sprintf("test_data_%d", i))
	}

	addr := "localhost:1222"
	for _, c := range config.NcpIO {
		if c.Type == "tcps" {
			addr = c.Params
		}
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	if strings.TrimSuffix(status, "\n") != msg1 {
		log.Panicln("Should", msg1)
	}

	status, err = reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	if strings.TrimSuffix(status, "\n") != msg2 {
		log.Panicln("Should", msg2)
	}

	conn.Write([]byte(msg2 + "\n"))
	conn.Write([]byte(msg1 + "\n"))

	if string(<-ncpio.O) != msg2 {
		log.Panicln("Should", msg2)
	}

	if string(<-ncpio.O) != msg1 {
		log.Panicln("Should", msg1)
	}

	log.Println("Socket Successfully")
}
