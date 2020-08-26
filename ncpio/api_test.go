package ncpio

import (
	"context"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	params := "233"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := make(chan []byte, 128)
	o := make(chan []byte, 128)
	go NewApi(params, i, o).Run(ctx)

	test_1 := `{"jsonrpc":"2.0","method":"dooropen","params":[]}`
	test_2 := `{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`

	i <- []byte(test_1)
	i <- []byte(test_2)

	I <- []byte(test_1)
	I <- []byte(test_2)

	if data := string(<-O); data != test_1 {
		t.Error(data)
	}
	if data := string(<-O); data != test_2 {
		t.Error(data)
	}

	if data := string(<-o); data != test_1 {
		t.Error(data)
	}
	if data := string(<-o); data != test_2 {
		t.Error(data)
	}
}

func TestApi2(t *testing.T) {
	// coverage
	// Need to wait for 'goroutine' to be closed by 'ctx' before ending 'Test' process
	time.Sleep(time.Millisecond)
}
