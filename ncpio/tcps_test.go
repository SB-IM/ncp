package ncpio

import (
	"bufio"
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

func TestTcps(t *testing.T) {
	addr := "localhost:1222"
	ncpios := NewNcpIOs([]Config{
		{
			Type: "api",
			IRules: []Rule{
				{`.*`, false},
			},
			ORules: []Rule{
				{`.*`, false},
			},
		},
		{
			Type:   "tcps",
			Params: addr,
			IRules: []Rule{
				{`.*`, false},
			},
			ORules: []Rule{
				{`.*`, false},
			},
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ncpios.Run(ctx)
	time.Sleep(3 * time.Millisecond)

	msg1 := "2333333333333"
	msg2 := "4555555555555"

	I <- []byte(msg1)
	I <- []byte(msg2)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("dial error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		t.Error(err)
	}
	if strings.TrimSuffix(status, "\n") != msg1 {
		t.Error("Should", msg1)
	}

	status, err = reader.ReadString('\n')
	if err != nil {
		t.Error(err)
	}
	if strings.TrimSuffix(status, "\n") != msg2 {
		t.Error("Should", msg2)
	}

	conn.Write([]byte(msg2 + "\n"))
	conn.Write([]byte(msg1 + "\n"))

	if string(<-O) != msg2 {
		t.Error("Should", msg2)
	}

	if string(<-O) != msg1 {
		t.Error("Should", msg1)
	}

	// Test TCPS Listener not closed
	cancel()
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	go ncpios.Run(ctx2)
	time.Sleep(time.Millisecond)
}
