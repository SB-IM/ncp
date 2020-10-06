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
		ncps[index] = NewNcpIO(index, &config)
	}
	return &NcpIOs{
		IO: ncps,
	}
}

func (n *NcpIOs) Run(ctx context.Context) {
	type ncpData struct {
		Name string
		Data []byte
	}

	hub := make(chan ncpData, ioChannelBuffering)

	for _, io := range n.IO {
		go io.Run(ctx)
		go func(ncpio *NcpIO) {
			for data := range ncpio.IO.O {
				if Filter(ncpio.ORules, data) {
					hub <- ncpData{
						Name: ncpio.Name,
						Data: data,
					}
				}
			}
		}(io)
	}

	for raw := range hub {
		for _, io := range n.IO {

			// Skip data loopback
			if io.Name == raw.Name {
				continue
			}

			// Skip No Match IRules
			if !Filter(io.IRules, raw.Data) {
				continue
			}

			select {
			case io.IO.I <- raw.Data:
			default:
			}
		}
	}
}
