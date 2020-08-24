package ncpio

import (
	"bufio"
	"context"
	"net"
	"testing"
)

const (
	testAddr = "localhost:8877"
)

func TestTcpc(t *testing.T) {
	sign := make(chan bool)

	listener, err := net.Listen("tcp", testAddr)
	if err != nil {
		t.Error(err)
	}

	test_1 := "23333333333333"
	test_2 := "4555555"

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
		}
		conn.Write([]byte(test_1 + "\n"))
		conn.Write([]byte(test_2 + "\n"))

		read := bufio.NewReader(conn)
		if data, err := readLine(read); err == nil {
			if string(data) != test_1 {
				t.Errorf("%s\n", data)
			}
		} else {
			t.Error(err)
		}

		if data, err := readLine(read); err == nil {
			if string(data) != test_2 {
				t.Errorf("%s\n", data)
			}
		} else {
			t.Error(err)
		}

		sign <- true
		conn.Close()

		sign <- true
		// Reconnect
		conn, err = listener.Accept()
		read = bufio.NewReader(conn)
		if data, err := readLine(read); err == nil {
			if string(data) != test_1 {
				t.Errorf("%s\n", data)
			}
		} else {
			t.Error(err)
		}

		conn.Write([]byte(test_1 + "\n"))
		sign <- true
	}()

	l := 128
	i := make(chan []byte, l)
	o := make(chan []byte, l)
	tcpc := NewTcpc(testAddr, i, o)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tcpc.Run(ctx)

	if msg := <-i; string(msg) != test_1 {
		t.Errorf("%s\n", msg)
	}
	if msg := <-i; string(msg) != test_2 {
		t.Errorf("%s\n", msg)
	}

	o <- []byte(test_1)
	o <- []byte(test_2)

	<-sign

	for i := 0; i <= l; i++ {
		o <- []byte(test_1)
	}
	<-sign
	if msg := <-i; string(msg) != test_1 {
		t.Errorf("%s\n", msg)
	}
	<-sign
}
