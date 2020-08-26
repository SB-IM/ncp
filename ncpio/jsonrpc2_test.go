package ncpio

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/SB-IM/jsonrpc2"
)

func TestJsonrpc2(t *testing.T) {
	params := "233"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := make(chan []byte)
	o := make(chan []byte)
	go NewJsonrpc2(params, i, o).Run(ctx)

	i <- []byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`)
	i <- []byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`)
	data := <-o

	j := &jsonrpc2.Jsonrpc{}
	json.Unmarshal(data, j)
	if !j.IsSuccess() {
		t.Errorf("%s\n", data)
	}

	if d, _ := j.Result.MarshalJSON(); string(d) != params {
		t.Errorf("%s\n", d)
	}

}

func TestJsonrpc2Error(t *testing.T) {
	params := `{"error": {"code": 0, "message": "xxxxx"}}`

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := make(chan []byte)
	o := make(chan []byte)
	go NewJsonrpc2(params, i, o).Run(ctx)

	i <- []byte("233")
	i <- []byte(params)
	i <- []byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`)
	data := <-o

	j := &jsonrpc2.Jsonrpc{}
	json.Unmarshal(data, j)
	if j.IsSuccess() {
		t.Errorf("%s\n", data)
	} else {
		if j.Error.Message != "xxxxx" {
			t.Error(j.Error.Message)
		}
	}
}
