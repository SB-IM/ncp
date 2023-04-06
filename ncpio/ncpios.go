package ncpio

import (
	"context"
)

type NcpIOs struct {
	IO      []*NcpIO
	Debuger Debuger
}

func NewNcpIOs(configs []Config) *NcpIOs {
	ncps := make([]*NcpIO, len(configs))
	for index, config := range configs {
		ncps[index] = NewNcpIO(index, &config)
	}
	return &NcpIOs{
		IO:      ncps,
		Debuger: NoDebuger{},
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
				n.Debuger.Printf("<%s> RECV: %s", ncpio.Name, data)
				if Filter(ncpio.ORules, data) {
					n.Debuger.Printf("<%s> CAST: %s", ncpio.Name, data)
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

			n.Debuger.Printf("<%s> SEND: %s", io.Name, raw.Data)

			select {
			case io.IO.I <- raw.Data:
			default:
			}
		}
	}
}
