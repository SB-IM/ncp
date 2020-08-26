package ncpio

import (
	"context"
	"time"
)

const (
	ioChannelBuffering = 128
	retryInterval      = 3 * time.Second
)

type NcpIO struct {
	IRules []Rule `json:"i_rules"`
	ORules []Rule `json:"o_rules"`
	Run    func(context.Context)
	IO     IO
}

type IO struct {
	I chan []byte
	O chan []byte
}

type Config struct {
	// tcps / tcpc / mqtt / history / logger / jsonrpc2 / build-in / api
	Type   string `json:"type" yaml:"type"`
	Params string `json:"params" yaml:"params"`
	IRules []Rule `json:"i_rules" yaml:"i_rules"`
	ORules []Rule `json:"o_rules" yaml:"o_rules"`
}

func NewNcpIO(config *Config) *NcpIO {
	i := make(chan []byte, ioChannelBuffering)
	o := make(chan []byte, ioChannelBuffering)

	run := func() func(context.Context) {
		switch config.Type {
		case "api":
			return NewApi(config.Params, i, o).Run
		case "jsonrpc2":
			return NewJsonrpc2(config.Params, i, o).Run
		default:
			return NewApi(config.Params, i, o).Run
		}
	}()

	return &NcpIO{
		IRules: config.IRules,
		ORules: config.ORules,
		Run:    run,
		IO: IO{
			I: i,
			O: o,
		},
	}
}
