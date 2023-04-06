package ncpio

import (
	"context"
	"testing"

	"github.com/sb-im/jsonrpc-lite"
)

func TestExecutor(t *testing.T) {
	params := "echo"
	echo := `"233vv"`

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := make(chan []byte)
	o := make(chan []byte)
	go NewExecutor(params, i, o).Run(ctx)

	i <- []byte(`{"jsonrpc":"2.0","method":"dooropen","params":["233"]}`)
	i <- []byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":["-n",` + echo + `]}`)
	data := <-o

	j := jsonrpc.ParseObject(data)
	if j.Type != jsonrpc.TypeSuccess {
		t.Errorf("%s\n", data)
	}

	if d, _ := j.Result.MarshalJSON(); string(d) != echo {
		t.Errorf("%s\n", d)
	}

}

func TestExecutorError(t *testing.T) {
	params := "sh"
	args := `"-n","exit 1"`

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i := make(chan []byte)
	o := make(chan []byte)
	go NewExecutor(params, i, o).Run(ctx)

	i <- []byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[` + args + `]}`)
	data := <-o

	j := jsonrpc.ParseObject(data)
	if j.Type == jsonrpc.TypeSuccess {
		t.Errorf("%s\n", data)
	} else {
		if j.Errors.Message != "exec error" {
			t.Error(j.Errors.Message)
		}
	}
}
