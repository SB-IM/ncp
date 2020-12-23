package ncpio

import (
	"testing"
)

func TestFilter(t *testing.T) {
	if Filter(
		[]Rule{
			{
				Regexp: "233",
			},
			{
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
			{
				Regexp: "233",
			},
			{
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
			{
				Regexp: "2(d",
			},
			{
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
			{
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

func TestFilterMultipleRules(t *testing.T) {
	if Filter(
		[]Rule{
			{
				Regexp: `.*"jsonrpc".*`,
				Invert: true,
			},
			{
				Regexp: "233",
			},
			{
				Regexp: `.*"method".*`,
				Invert: true,
			},
		},
		[]byte(`{"jsonrpc":"2.0","method":"dooropen","params":[]}`),
	) {
		t.Error("Not Match")
	}
}

func TestFilterMultipleRules2(t *testing.T) {
	if Filter(
		[]Rule{
			{
				Regexp: `.*"method": ?"webrtc".*`,
				Invert: true,
			},
			{
				Regexp: `.*`,
				Invert: true,
			},
		},
		[]byte(`{"jsonrpc":"2.0","id":"test.0-0","method":"webrtc","params":[]}`),
	) {
		t.Error("Not Match")
	}
}
