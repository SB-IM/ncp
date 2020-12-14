package util

import (
	"testing"
)

func TestBinaryPrefix(t *testing.T) {
	if r, _ := BinaryPrefix("128M"); r != 134217728 {
		t.Error(r)
	}

	if r, _ := BinaryPrefix("128"); r != 128 {
		t.Error(r)
	}

	if r, err := BinaryPrefix(""); r != 0 && err != nil {
		t.Error(r)
	}

	if r, err := BinaryPrefix("MMMM"); r != 0 && err != nil {
		t.Error(r)
	}
}
