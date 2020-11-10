package ncpio

import (
	"context"
	"io"
	"log"
	"os"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type Logger struct {
	Log *log.Logger
	I   <-chan []byte
	O   chan<- []byte
}

func NewLogger(params string, i <-chan []byte, o chan<- []byte) *Logger {
	prefix := ""

	var out io.Writer
	out, err := rotatelogs.New(
		params+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(params),
		// Max 128M
		rotatelogs.WithRotationSize(128*1024*1024),
		rotatelogs.WithRotationCount(8),
	)

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
