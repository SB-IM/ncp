package ncpio

import (
	"testing"
)

func TestGet(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
		ORules: []Rule{
			Rule{`.*`, true},
			Rule{`23`, false},
		},
	}

	tcpc := NewTcpc(ncpio.Params)
	ncpio.Conn = tcpc

	tcpc.output <- []byte("23333333333333\n")
	if data, err := ncpio.Get(); err != nil {
		t.Errorf("%s\n", data)
	}

	tcpc.output <- []byte("344444444\n")
	if data, err := ncpio.Get(); err == nil {
		t.Errorf("%s\n", data)
		t.Errorf("%s\n", err)
	}

	// No Data
	if data, err := ncpio.Get(); err == nil {
		t.Errorf("%s\n", data)
		t.Errorf("%s\n", err)
	}
}

func TestGetNot(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
		IRules: []Rule{
			Rule{`.*`, true},
			Rule{`.*`, false},
		},
	}

	tcpc := NewTcpc(ncpio.Params)
	ncpio.Conn = tcpc

	tcpc.output <- []byte("23333333333333\n")

	if data, err := ncpio.Get(); err == nil {
		t.Errorf("%s\n", data)
		t.Errorf("%s\n", err)
	}
}

func TestGetErr(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
		ORules: []Rule{
			Rule{`*`, true},
		},
	}

	tcpc := NewTcpc(ncpio.Params)
	ncpio.Conn = tcpc

	tcpc.output <- []byte("23333333333333\n")
	if data, err := ncpio.Get(); err == nil {
		t.Errorf("%s\n", data)
		t.Errorf("%s\n", err)
	}
}

func TestPut(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
		IRules: []Rule{
			Rule{`.*`, true},
			Rule{`23`, false},
		},
	}

	tcpc := NewTcpc(ncpio.Params)
	ncpio.Conn = tcpc

	data := []byte("23333333333333")

	if err := ncpio.Put([]byte(data)); err != nil {
		t.Errorf("%s\n", err)
	}

	if d, err := tcpc.Get(); string(data) == string(d) {
		t.Errorf("%s\n", d)
		t.Errorf("%s\n", err)
	}

	if d, err := tcpc.Get(); string(data) == string(d) {
		t.Errorf("%s\n", d)
		t.Errorf("%s\n", err)
	}
}

func TestPutErr(t *testing.T) {
	ncpio := &NcpIO{
		Type:   "tcpc",
		Params: "233",
	}

	tcpc := NewTcpc(ncpio.Params)
	ncpio.Conn = tcpc

	data := []byte("23333333333333")

	if err := ncpio.Put([]byte(data)); err != nil {
		t.Errorf("%s\n", err)
	}

	if d, err := tcpc.Get(); string(data) == string(d) {
		t.Errorf("%s\n", d)
		t.Errorf("%s\n", err)
	}

	if d, err := tcpc.Get(); string(data) == string(d) {
		t.Errorf("%s\n", d)
		t.Errorf("%s\n", err)
	}
}
