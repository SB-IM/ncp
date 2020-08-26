package ncpio

import (
	"testing"
)

func TestFilter(t *testing.T) {
	if Filter(
		[]Rule{
			Rule{
				Regexp: "233",
			},
			Rule{
				Regexp: `.*"result".*`,
			},
		},
		[]byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`),
	) {
		t.Error("Not Match")
	}
}

func TestFilterMatch(t *testing.T) {
	if !Filter(
		[]Rule{
			Rule{
				Regexp: "233",
			},
			Rule{
				Regexp: `.*"method".*`,
			},
		},
		[]byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`),
	) {
		t.Error("Match")
	}
}

func TestFilterRegexpErrorAfterMatch(t *testing.T) {
	if !Filter(
		[]Rule{
			Rule{
				Regexp: "2(d",
			},
			Rule{
				Regexp: `.*"method".*`,
			},
		},
		[]byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`),
	) {
		t.Error("Match")
	}
}

func TestFilterMatchInvert(t *testing.T) {
	if Filter(
		[]Rule{
			Rule{
				Regexp: `.*"method".*`,
				Invert: true,
			},
		},
		[]byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`),
	) {
		t.Error("Not Match")
	}
}

func TestFilterNoRule(t *testing.T) {
	if Filter(
		[]Rule{},
		[]byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`),
	) {
		t.Error("Not Match")
	}
}
