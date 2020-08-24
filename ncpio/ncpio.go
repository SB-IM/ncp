package ncpio

import (
	"context"
	"errors"
	"regexp"
	"time"
)

const (
	retryInterval = 3 * time.Second
)

type NcpIO struct {
	// tcps / tcpc / mqtt / history / logger / jsonrpc2 / build-in
	Type   string `json:"type"`
	Params string `json:"params"`
	IRules []Rule `json:"i_rules"`
	ORules []Rule `json:"o_rules"`
	Conn   IOServer
	I      chan<- []byte
	O      <-chan []byte
	ConnI  chan<- []byte
	ConnO  <-chan []byte
}

type Rule struct {
	Regexp string `json:"regexp"`
	Invert bool   `json:"invert"`
}

type IOServer interface {
	Run(context.Context)
	//Run(context.Context) (chan<- []byte, chan-> []byte)
	Get() ([]byte, error)
	Put([]byte) error
	//I() chan<- []byte
}

func (n *NcpIO) Run(ctx context.Context) {
	switch n.Type {
	case "tcpc":
		n.Conn = NewTcpc(n.Params)
		go n.Conn.Run(ctx)
	case "jsonrpc2":
		n.Conn = NewJsonrpc2(n.Params)
		go n.Conn.Run(ctx)
	}
}

func (n *NcpIO) Get() ([]byte, error) {
	data, err := n.Conn.Get()
	if err != nil {
		return data, err
	}

	for _, rule := range n.ORules {
		matched, err := regexp.Match(rule.Regexp, data)
		if err != nil {
			return data, err
		}

		// Invert result
		if rule.Invert {
			matched = !matched
		}

		if matched {
			return data, err
		}
	}
	return data, errors.New("Not Match")
}

func (n *NcpIO) Put(data []byte) error {
	for _, rule := range n.IRules {
		matched, err := regexp.Match(rule.Regexp, data)
		if err != nil {
			return err
		}

		// Invert result
		if rule.Invert {
			matched = !matched
		}

		if matched {
			return n.Conn.Put(data)
		}
	}
	return errors.New("Not Match")
}
