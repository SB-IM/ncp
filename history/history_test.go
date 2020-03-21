package history

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	archive := New(100)

	data := make(map[string][]byte)
	data["test"] = []byte("2333333333333")

	archive.Add(data)

	if bytes.Compare(archive.GetLatest("test233"), data["test"]) == 0 {
		t.Errorf(string(data["test"]))
	}

	if bytes.Compare(archive.GetLatest("test"), data["test"]) != 0 {
		t.Errorf(string(data["test"]))
	}

	data2 := make(map[string][]byte)
	data2["test"] = []byte("4555555555555")
	archive.Add(data2)
	if bytes.Compare(archive.GetLatest("test"), data2["test"]) != 0 {
		t.Errorf(string(data2["test"]))
	}
}

func Test_GetLatestHistorys(t *testing.T) {
	tag := "test"
	archive := New(100)

	data := make(map[string][]byte)
	data[tag] = []byte("2333333333333")

	for i := 0; i <= 10; i++ {
		archive.Add(data)
	}

	time.Sleep(10 * time.Millisecond)

	data2 := make(map[string][]byte)
	data2[tag] = []byte("4555555555555")
	count := 3
	for i := 0; i <= count; i++ {
		archive.Add(data2)
	}

	if historys := archive.GetLatestHistorys(tag, "5ms"); len(historys) != count {
		t.Errorf("Total: %d, Should match: %d, Actual match: %d\n", len(archive.Tag[tag]), count, len(historys))
		for _, item := range historys {
			t.Errorf("Time: %d, Data: %s\n", item.Time.UnixNano(), item.Data)
		}
	}

}

func Test_GetLatestHistorys_nil(t *testing.T) {
	tag := "test"
	archive := New(100)

	if historys := archive.GetLatestHistorys(tag, "5ms"); len(historys) != 0 {
		t.Errorf("Should no Data\n")
		for _, item := range historys {
			t.Errorf("Time: %d, Data: %s\n", item.Time.UnixNano(), item.Data)
		}
	}
}

func Test_GetLatestHistorys_one(t *testing.T) {
	tag := "test"
	archive := New(100)

	data := make(map[string][]byte)
	data[tag] = []byte("2333333333333")

	archive.Add(data)

	if historys := archive.GetLatestHistorys(tag, "5ms"); len(historys) != 0 {
		t.Errorf("Should no Data\n")
		for _, item := range historys {
			t.Errorf("Time: %d, Data: %s\n", item.Time.UnixNano(), item.Data)
		}
	}
}

func Test_GetLatestHistorys_side(t *testing.T) {
	tag := "test"
	archive := New(100)

	data := make(map[string][]byte)
	data[tag] = []byte("2333333333333")

	archive.Add(data)

	time.Sleep(10 * time.Millisecond)

	data2 := make(map[string][]byte)
	data2[tag] = []byte("4555555555555")
	count := 3
	for i := 0; i <= count; i++ {
		archive.Add(data2)
	}

	if historys := archive.GetLatestHistorys(tag, "5ms"); len(historys) != count {
		t.Errorf("Total: %d, Should match: %d, Actual match: %d\n", len(archive.Tag[tag]), count, len(historys))
		for _, item := range historys {
			t.Errorf("Time: %d, Data: %s\n", item.Time.UnixNano(), item.Data)
		}
	}

	time.Sleep(10 * time.Millisecond)
	archive.Add(data)

	if historys := archive.GetLatestHistorys(tag, "5ms"); len(historys) != 1 {
		t.Errorf("Total: %d, Should match: %d, Actual match: %d\n", len(archive.Tag[tag]), count, len(historys))
		for _, item := range historys {
			t.Errorf("Time: %d, Data: %s\n", item.Time.UnixNano(), item.Data)
		}
	}

}

func Benchmark_GetLatestHistorys_1000000(b *testing.B) {
	tag := "test"
	archive := New(100)

	data := make(map[string][]byte)
	data[tag] = []byte("2333333333333")

	for i := 0; i <= 1000000; i++ {
		archive.Add(data)
	}

	for n := 0; n < b.N; n++ {
		archive.Add(data)
		archive.GetLatestHistorys(tag, "5ms")
	}
}

func Test_FilterAdd(t *testing.T) {
	tag := "test"
	archive := New(100)

	if !archive.FilterAdd(tag, []byte("2333333333333")) {
		t.Errorf("is True")
	}

	if archive.FilterAdd(tag, []byte("2333333333333")) {
		t.Errorf("is False")
	}

	if !archive.FilterAdd(tag, []byte("45555555")) {
		t.Errorf("is True")
	}

	if !archive.FilterAdd(tag, []byte("2333333333333")) {
		t.Errorf("is True")
	}
}

func Test_FilterAdd_Max(t *testing.T) {
	tag := "test"
	max := 100
	archive := New(max)

	for n := 0; n < max+10; n++ {
		if !archive.FilterAdd(tag, []byte(strconv.Itoa(n))) {
			t.Errorf("is True")
		}
	}

	if l := len(archive.Tag[tag]); l != max {
		t.Errorf("len: %d", l)
	}
}
