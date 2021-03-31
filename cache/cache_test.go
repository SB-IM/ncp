package cache

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	lRUCache := NewLRUCache(2)
	lRUCache.Put("1", "1")
	lRUCache.Put("2", "2")
	if lRUCache.Get("1") != "1" {
		t.Error("1")
	}
	lRUCache.Put("3", "3")

	if lRUCache.Get("2") == "2" {
		t.Error("2")
	}
	lRUCache.Put("4", "4")

	if lRUCache.Get("1") == "1" {
		t.Error("1")
	}
	if lRUCache.Get("3") != "3" {
		t.Error("3")
	}
	if lRUCache.Get("4") != "4" {
		t.Error("4")
	}
}
