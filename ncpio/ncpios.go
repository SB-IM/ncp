package ncpio

import (
	"context"
)

type NcpIOs struct {
	IO []*NcpIO
}

func NewNcpIOs(configs []Config) *NcpIOs {
	ncps := make([]*NcpIO, len(configs))
	for index, config := range configs {
		ncps[index] = NewNcpIO(&config)
	}
	return &NcpIOs{
		IO: ncps,
	}
}

func (n *NcpIOs) Run(ctx context.Context) {
	hub := make(chan []byte, ioChannelBuffering)

	for _, io := range n.IO {
		go io.Run(ctx)
		go func(ncpio *NcpIO) {
			for data := range ncpio.IO.O {
				if Filter(ncpio.ORules, data) {
					hub <- data
				}
			}
		}(io)
	}

	for data := range hub {
		for _, io := range n.IO {
			if !Filter(io.IRules, data) {
				continue
			}
			select {
			case io.IO.I <- data:
			default:
			}
		}
	}
}
