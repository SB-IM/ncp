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
		Config{
			Type: "api",
			IRules: []Rule{
				Rule{`.*`, false},
			},
			ORules: []Rule{
				Rule{`.*`, false},
			},
		},
		Config{
			Type:   "tcps",
			Params: addr,
			IRules: []Rule{
				Rule{`.*`, false},
			},
			ORules: []Rule{
				Rule{`.*`, false},
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
}
