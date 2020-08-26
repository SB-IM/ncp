package ncpio

import (
	"context"
)

const (
	apiChannelBuffering = 128
)

var (
	I chan []byte
	O chan []byte
)

type Api struct {
	I <-chan []byte
	O chan<- []byte
}

func NewApi(params string, i <-chan []byte, o chan<- []byte) *Api {
	I = make(chan []byte, apiChannelBuffering)
	O = make(chan []byte, apiChannelBuffering)
	return &Api{
		I: i,
		O: o,
	}
}

func (t *Api) Run(ctx context.Context) {
	for {
		select {
		case data := <-t.I:
			O <- data
		case data := <-I:
			t.O <- data
		case <-ctx.Done():
			return
		}
	}
}
