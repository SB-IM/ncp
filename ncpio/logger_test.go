package ncpio

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	TestLoggerPath = "/tmp/ncp_test/"
)

func TestLogger(t *testing.T) {
	params := "file://"+TestLoggerPath+"test.log?size=128M&count=8&prefix=TEST"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := make(chan []byte, 128)
	o := make(chan []byte, 128)
	go NewLogger(params, i, o).Run(ctx)

	l := 10000
	for n := 0; n < l; n++ {
		i <- []byte("TTTTTTTTTTTTTTTTTT")
	}

	// Need Wait log Write Completed
	time.Sleep(3 * time.Millisecond)

	file, err := os.Open(TestLoggerPath + "test.log")
	if err != nil {
		t.Error(err)
	}

	count, err := lineCounter(file)
	if err != nil {
		t.Error(err)
	}

	if count != l {
		t.Error("Line Counter not", count)
	}

	os.RemoveAll(TestLoggerPath)
}

func lineCounter(r io.Reader) (int, error) {
    buf := make([]byte, 32*1024)
    count := 0
    lineSep := []byte{'\n'}

    for {
        c, err := r.Read(buf)
        count += bytes.Count(buf[:c], lineSep)

        switch {
        case err == io.EOF:
            return count, nil

        case err != nil:
            return count, err
        }
    }
}

func TestRotateLogger(t *testing.T) {
	params := "file://"+TestLoggerPath+"rotate.log?size=128&count=5&prefix=TEST"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := make(chan []byte, 128)
	o := make(chan []byte, 128)
	go NewLogger(params, i, o).Run(ctx)

	for n := 0; n < 10000; n++ {
		i <- []byte("TTTTTTTTTTTTTTTTTT")
	}

	file, err := ioutil.ReadDir(TestLoggerPath)
	if err != nil {
		t.Error(err)
	}

	if len(file) < 3 {
		t.Error("Log File not rotate")
	}

	os.RemoveAll(TestLoggerPath)
}
