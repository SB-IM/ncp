package ncpio

import (
	"context"
	"encoding/json"

	"github.com/sb-im/jsonrpc-lite"
	logger "log"
)

type Jsonrpc2 struct {
	jsonrpc *jsonrpc.Jsonrpc
	I       <-chan []byte
	O       chan<- []byte
}

func NewJsonrpc2(params string, i <-chan []byte, o chan<- []byte) *Jsonrpc2 {
	jsonrpc_res := &jsonrpc.Jsonrpc{}

	err := json.Unmarshal([]byte(params), jsonrpc_res)
	if err != nil {
		raw := json.RawMessage([]byte(params))
		jsonrpc_res.Result = &raw
	}

	return &Jsonrpc2{
		jsonrpc: jsonrpc_res,
		I:       i,
		O:       o,
	}
}

func (t *Jsonrpc2) Run(ctx context.Context) {
	t.simulation(ctx)
}

func (t *Jsonrpc2) simulation(ctx context.Context) {
	for {
		select {
		case raw := <-t.I:
			data, err := t.rpcCall(raw)
			if err != nil {
				logger.Println(string(data), err)
				continue
			}
			if len(data) != 0 {
				t.O <- data
			}
		case <-ctx.Done():
			return
		}
	}
}

func (t *Jsonrpc2) rpcCall(req []byte) ([]byte, error) {
	jsonrpc_req, err := jsonrpc.Parse(req)

	if err != nil {
		return []byte{}, err
	}

	if jsonrpc_req.Type == jsonrpc.TypeInvalid {
		return []byte{}, nil
	}

	if jsonrpc_req.Type == jsonrpc.TypeNotify {
		return []byte{}, nil
	}

	t.jsonrpc.ID = jsonrpc_req.ID
	return json.Marshal(t.jsonrpc)
}
