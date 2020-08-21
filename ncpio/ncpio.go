package ncpio

import (
	"context"
	"errors"
	"regexp"
)

type NcpIO struct {
	// tcps / tcpc / mqtt / history / logger / jsonrpc2 / build-in
	Type   string `json:"type"`
	Params string `json:"params"`
	IRules []Rule `json:"i_rules"`
	ORules []Rule `json:"o_rules"`
	Conn   IOServer
}

type Rule struct {
	Regexp string `json:"regexp"`
	Invert bool   `json:"invert"`
}

type IOServer interface {
	Run(context.Context)
	Get() ([]byte, error)
	Put([]byte) error
}

func (n *NcpIO) Run(ctx context.Context) {
	if n.Type == "tcpc" {
		go NewTcpc(n.Params).Run(ctx)
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
			return err
		}
	}
	return n.Conn.Put(data)
}
