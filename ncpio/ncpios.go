package ncpio

import (
	"context"
)

type NcpIOs struct {
	IO []*NcpIO
}

func (n *NcpIOs) Run(ctx context.Context) {
	for _, io := range n.IO {
		io.Run(ctx)
	}

	for _, io := range n.IO {
		io.Get()
	}
}
