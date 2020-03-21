package history

import (
	//"fmt"
	"bytes"
	"time"
)

type History struct {
	Time time.Time
	Data []byte
}

type Archive struct {
	Tag map[string][]History
	Len int
}

func New(max int) *Archive {
	archive := &Archive{
		Tag: make(map[string][]History),
		Len: max,
	}
	return archive
}

func (this *Archive) Add(dataMap map[string][]byte) {
	for key, value := range dataMap {
		history := History{
			Time: time.Now(),
			Data: value,
		}
		this.Tag[key] = append(this.Tag[key], history)
	}
}

func (this *Archive) GetLatest(name string) []byte {
	if this.Tag[name] == nil {
		return []byte("")
	}

	return this.Tag[name][len(this.Tag[name])-1].Data
}

func (this *Archive) FilterAdd(name string, data []byte) bool {
	if bytes.Compare(this.GetLatest(name), data) != 0 {
		history := History{
			Time: time.Now(),
			Data: data,
		}

		if l := len(this.Tag[name]); l >= this.Len {
			this.Tag[name] = this.Tag[name][l-this.Len : l-1]
		}
		this.Tag[name] = append(this.Tag[name], history)
		return true
	}

	return false
}

func (this *Archive) GetLatestHistorys(name, strDuration string) []History {
	if len(this.Tag[name]) == 0 {
		return []History{}
	}
	duration, _ := time.ParseDuration("-" + strDuration)
	sinceTime := time.Now().Add(duration)

	var getIncludeIndex func(int, int) int
	getIncludeIndex = func(i, num int) int {
		n := num / 2
		//fmt.Printf("i: %d, num: %d, n: %d, r: %d === ", i, num, n, r)
		if n == 0 {
			//fmt.Println("getIncludeIndex: SUCCESS")
			return i
		} else if sinceTime.Before(this.Tag[name][i+n].Time) {
			i = getIncludeIndex(i, n)
		} else {
			i = getIncludeIndex(i+n, n)
		}
		//fmt.Println("getIncludeIndex: END")
		return i
	}

	index := getIncludeIndex(0, len(this.Tag[name])-1)
	//fmt.Println(index)
	if len(this.Tag[name]) != 1 {
		index++
	}
	return this.Tag[name][index : len(this.Tag[name])-1]
}
