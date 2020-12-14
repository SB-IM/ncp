package ncpio

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"

	"sb.im/ncp/util"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type Logger struct {
	Log *log.Logger
	I   <-chan []byte
	O   chan<- []byte
}

func NewLogger(params string, i <-chan []byte, o chan<- []byte) *Logger {
	prefix := "NCPIO"

	// Default Max 128M
	var size int64 = 128 * 1024 * 1024
	// Default Count 0
	var count int

	if u, err := url.Parse(params); err != nil {
		params = ""
	} else {
		params = u.Path
		q := u.Query()
		if s, e := util.BinaryPrefix(q.Get("size")); e == nil {
			size = s
		}
		count, _ = strconv.Atoi(q.Get("count"))
		if s := q.Get("prefix"); s != "" {
			prefix = s
		}
	}

	var err error
	var out io.Writer
	if params == "" {
		out = os.Stdout
		prefix = "[DEV] NCPIO"
	} else {
		out, err = rotatelogs.New(
			params+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(params),
			rotatelogs.WithRotationSize(size),
			rotatelogs.WithRotationCount(uint(count)),
		)
		if err != nil {
			out = os.Stdout
			prefix = "FILE ERROR"
		}
	}

	return &Logger{
		Log: log.New(out, fmt.Sprintf("[%s]: ", prefix), log.LstdFlags),
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
