package ncpio

import (
	"encoding/json"
	"testing"

	"github.com/SB-IM/jsonrpc2"
)

func TestJsonrpc2(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
		ORules: []Rule{
			Rule{`.*`, true},
			Rule{`23`, false},
		},
	}

	jsonrpc := NewJsonrpc2(ncpio.Params)
	err := jsonrpc.Put([]byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`))
	if err != nil {
		t.Error(err)
	}

	data, err := jsonrpc.Get()
	if err != nil {
		t.Error(err)
	}

	j := &jsonrpc2.Jsonrpc{}
	json.Unmarshal(data, j)
	if !j.IsSuccess() {
		t.Errorf("%s\n", data)
	}
}

func TestJsonrpc2Notify(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
		ORules: []Rule{
			Rule{`.*`, true},
			Rule{`23`, false},
		},
	}

	raw := `{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`

	jsonrpc := NewJsonrpc2(ncpio.Params)

	sum := 0
	for i := 0; i <= 127; i++ {
		sum += i

		err := jsonrpc.Put([]byte(raw))
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i <= 25; i++ {
		sum += i

		err := jsonrpc.Put([]byte(raw))
		if err == nil {
			t.Error(err)
		}
	}
}

func TestJsonrpc2Result(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: `{"result":{"x": "xxxxx"}}`,
		ORules: []Rule{
			Rule{`.*`, true},
			Rule{`23`, false},
		},
	}

	jsonrpc := NewJsonrpc2(ncpio.Params)
	err := jsonrpc.Put([]byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`))
	if err != nil {
		t.Error(err)
	}

	data, err := jsonrpc.Get()
	if err != nil {
		t.Error(err)
	}

	j := &jsonrpc2.Jsonrpc{}
	json.Unmarshal(data, j)
	if !j.IsSuccess() {
		t.Errorf("%s\n", data)
	}
}

func TestJsonrpc2Error(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: `{"error": {"code": 0, "message": "xxxxx"}}`,
		ORules: []Rule{
			Rule{`.*`, true},
			Rule{`23`, false},
		},
	}

	jsonrpc := NewJsonrpc2(ncpio.Params)
	err := jsonrpc.Put([]byte(`{"jsonrpc":"2.0","id":"sdwc.1-1553321035000","method":"dooropen","params":[]}`))
	if err != nil {
		t.Error(err)
	}

	data, err := jsonrpc.Get()
	if err != nil {
		t.Error(err)
	}

	j := &jsonrpc2.Jsonrpc{}
	json.Unmarshal(data, j)
	if j.IsSuccess() {
		t.Errorf("%s\n", data)
	}
}
