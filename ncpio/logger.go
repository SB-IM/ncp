package ncpio

import (
	"context"
	"log"
	"os"
)

type Logger struct {
	Log *log.Logger
	I   <-chan []byte
	O   chan<- []byte
}

func NewLogger(params string, i <-chan []byte, o chan<- []byte) *Logger {
	prefix := ""
	out, err := os.OpenFile(params, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		out = os.Stdout
		prefix = "NCPIO: "
	}
	return &Logger{
		Log: log.New(out, prefix, log.LstdFlags),
		I:   i,
		O:   o,
	}
}

func (t *Logger) Run(ctx context.Context) {
	for {
		select {
		case data := <-t.I:
			t.Log.Printf("%s\n", data)
		case <-ctx.Done():
			return
		}
	}
}
