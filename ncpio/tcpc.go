package ncpio

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	logger "log"
)

type Tcpc struct {
	conn   net.Conn
	addr   string
	buf    []byte
	mu_s   sync.Mutex
	mu_r   sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	I      chan<- []byte
	O      <-chan []byte
}

func NewTcpc(addr string, i chan<- []byte, o <-chan []byte) *Tcpc {
	return &Tcpc{
		addr: addr,
		I:    i,
		O:    o,
	}
}

func (t *Tcpc) Run(ctx context.Context) {
	for {
		conn, err := net.Dial("tcp", t.addr)
		if err != nil {
			time.Sleep(retryInterval)
		} else {
			logger.Println("New Connect")

			t.ctx, t.cancel = context.WithCancel(context.Background())
			go t.send(conn)
			go t.recv(conn)
			select {
			case <-t.ctx.Done():
				conn.Close()
				logger.Println("Connect err try reconnect")
			}
		}
	}
}

func (t *Tcpc) send(conn net.Conn) {
	t.mu_s.Lock()
	defer t.mu_s.Unlock()

	if len(t.buf) != 0 {
		_, err := conn.Write(append(t.buf, "\n"...))

		if err != nil {
			t.cancel()
			return
		}
		t.buf = []byte{}
	}

	for {
		select {
		case data := <-t.O:
			_, err := conn.Write(append(data, "\n"...))

			if err != nil {
				t.buf = data
				t.cancel()
				return
			}
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *Tcpc) recv(conn net.Conn) {
	t.mu_r.Lock()
	defer t.mu_r.Unlock()
	read := bufio.NewReader(conn)
	for {

		// readLine()
		data, err := readLine(read)
		if err != nil {
			t.cancel()
			return
		}
		t.I <- data
	}
}

// This function mainly solves the case where the number of bytes in a single line is greater than 4096
func readLine(reader *bufio.Reader) ([]byte, error) {
	var buffer bytes.Buffer
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return buffer.Bytes(), err
		}
		buffer.Write(line)
		if !isPrefix {
			break
		}
	}
	return buffer.Bytes(), nil
}
