package ncpio

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/SB-IM/jsonrpc2"
)

type Jsonrpc2 struct {
	jsonrpc *jsonrpc2.WireResponse
	output  chan []byte
}

func NewJsonrpc2(params string) *Jsonrpc2 {
	jsonrpc_res := &jsonrpc2.WireResponse{}

	err := json.Unmarshal([]byte(params), jsonrpc_res)
	if err != nil {
		raw := json.RawMessage([]byte(params))
		jsonrpc_res.Result = &raw
	}

	return &Jsonrpc2{
		jsonrpc: jsonrpc_res,
		output:  make(chan []byte, 128),
	}
}

func (t *Jsonrpc2) Run(ctx context.Context) {}

func (t *Jsonrpc2) Get() ([]byte, error) {
	select {
	case data := <-t.output:
		return data, nil
	default:
		return []byte{}, errors.New("Not get")
	}
}

func (t *Jsonrpc2) Put(req []byte) error {
	jsonrpc_req := jsonrpc2.WireRequest{}
	err := json.Unmarshal(req, &jsonrpc_req)
	if err != nil {
		return err
	}

	if jsonrpc_req.IsNotify() {
		return nil
	}

	t.jsonrpc.ID = jsonrpc_req.ID

	res, err := json.Marshal(t.jsonrpc)
	if err != nil {
		return err
	}

	select {
	case t.output <- res:
		return nil
	default:
		return errors.New("Not put")
	}
}
