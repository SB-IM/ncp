package ncpio

import (
	"context"
	"testing"
	"time"
)

func TestNcpIOs(t *testing.T) {
	configs := []Config{
		Config{
			Type:   "api",
			Params: "233",
			IRules: []Rule{
				Rule{`.*"result".*`, false},
			},
			ORules: []Rule{
				Rule{`.*"result".*`, true},
			},
		},
		Config{
			Type:   "jsonrpc2",
			Params: "233",
			IRules: []Rule{
				//Rule{`.*"result".*`, true},
				Rule{`.*`, false},
			},
			ORules: []Rule{
				Rule{`.*"result".*`, false},
			},
		},
	}

	ncpios := NewNcpIOs(configs)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ncpios.Run(ctx)
	time.Sleep(3 * time.Millisecond)

	test_1 := `{"jsonrpc":"2.0","method":"dooropen","params":[]}`
	test_2 := `{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`

	I <- []byte(test_2)
	I <- []byte(test_1)
	I <- []byte(test_2)

	if data := string(<-O); data == test_2 {
		t.Error(data)
	}
	if data := string(<-O); data == test_2 {
		t.Error(data)
	}
}
