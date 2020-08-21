package ncpio

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"net"
	"time"
)

type Tcpc struct {
	Conn   net.Conn
	addr   string
	input  chan []byte
	output chan []byte
}

func NewTcpc(addr string) *Tcpc {
	return &Tcpc{
		addr:   addr,
		input:  make(chan []byte, 128),
		output: make(chan []byte, 128),
	}
}

func (t *Tcpc) Run(ctx context.Context) {
	for {
		conn, err := net.Dial("tcp", t.addr)
		if err != nil {
			time.Sleep(3 * time.Second)
		} else {
			t.recv(conn, t.output)
			//this.logger.Println("Connect err try reconnect")
		}
	}
}

func (t *Tcpc) Get() ([]byte, error) {
	select {
	case data := <-t.output:
		return data, nil
	default:
		return []byte{}, errors.New("Not get")
	}
}

func (t *Tcpc) Put(data []byte) error {
	select {
	case t.input <- data:
		return nil
	default:
		return errors.New("Not put")
	}
}

func (t *Tcpc) recv(conn net.Conn, ch chan []byte) {
	read := bufio.NewReader(conn)
	for {

		// readLine()
		raw, err := readLine(read)
		if err != nil {
			conn.Close()
			return
		}
		ch <- raw
	}
}

func (t *Tcpc) send(conn net.Conn, ch chan []byte) {
	for msg := range ch {
		_, err := conn.Write(append(msg, "\n"...))

		if err != nil {
			ch <- msg
			break
		}
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
