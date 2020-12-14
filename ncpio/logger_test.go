package ncpio

import (
	"context"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	params := "file:///tmp/ncp/test.log?size=128M&count=8&prefix=SB"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := make(chan []byte, 128)
	o := make(chan []byte, 128)
	go NewLogger(params, i, o).Run(ctx)

	for n := 0; n < 10000; n++ {
		i <- []byte("TTTTTTTTTTTTTTTTTT")
	}

	time.Sleep(3 * time.Millisecond)
}
