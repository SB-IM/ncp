package ncpio

import (
	"context"
	"fmt"
	"time"
)

type NcpIOs struct {
	IO []*NcpIO
}

func (n *NcpIOs) Run(ctx context.Context) {
	for _, io := range n.IO {
		io.Run(ctx)
	}
	time.Sleep(3 * time.Second)

	for {
		for _, g_io := range n.IO {
			data, err := g_io.Get()
			if err != nil {
				continue
			}
			fmt.Printf("%s\n", data)

			for _, p_io := range n.IO {
				err := p_io.Put(data)
				if err != nil {
					fmt.Println(err)
				}
			}

		}
		time.Sleep(1 * time.Millisecond)
	}
}
