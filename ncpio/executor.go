package ncpio

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/sb-im/jsonrpc-lite"
)

type Executor struct {
	prefix string
	ctx    context.Context
	I      <-chan []byte
	O      chan<- []byte
}

func NewExecutor(params string, i <-chan []byte, o chan<- []byte) *Executor {
	return &Executor{
		prefix: params,
		I:      i,
		O:      o,
	}
}

func (t *Executor) Run(ctx context.Context) {
	t.ctx = ctx
	for {
		select {
		case raw := <-t.I:
			data, err := t.rpcCall(raw)
			if err != nil {
				log.Println(string(data), err)
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

func (t *Executor) doRun(args []string) ([]byte, error) {
	return exec.CommandContext(t.ctx, t.prefix, args...).CombinedOutput()
}

func (t *Executor) rpcCall(req []byte) ([]byte, error) {
	rpc, err := jsonrpc.Parse(req)
	var args []string

	if err := json.Unmarshal(*rpc.Params, &args); err != nil {
		return []byte{}, err
	}

	if rpc.Type == jsonrpc.TypeInvalid {
		return []byte{}, nil
	}

	raw, err := t.doRun(args)
	if rpc.Type == jsonrpc.TypeNotify {
		return []byte{}, nil
	}

	data, _ := json.Marshal(string(raw))
	d := json.RawMessage(data)

	if err != nil {
		return jsonrpc.NewError(rpc.ID, 0, "exec error", json.RawMessage(err.Error())).ToJSON()
	}
	return jsonrpc.NewSuccess(rpc.ID, d).ToJSON()
}
