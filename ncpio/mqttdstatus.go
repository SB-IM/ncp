package ncpio

import (
	"strconv"
	"time"
)

type StaticStatus struct {
	Link_id string `json:"link_id"`
	Type    string `json:"type"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	Alt     string `json:"alt"`
}

type NetworkStatus struct {
	// (0 ~ 100)%
	Loss int `json:"loss"`
	// AvgRtt (ms)
	Delay int `json:"delay"`
}

// {"code":0,"msg":"online","timestamp":"1607566337","status":{"link_id":3,"position_ok":true,"lat":"22.831878298","lng":"113.514688221","alt":"80.0001"}}
type NodeStatus struct {
	Code      int           `json:"code"`
	Msg       string        `json:"msg"`
	Timestamp string        `json:"timestamp"`
	Status    *StaticStatus `json:"status"`
}

func (t *NodeStatus) SetOnline(str string) *NodeStatus {
	statusMap := map[string]int{
		"online":   0,
		"offline":  1,
		"neterror": 2,
	}

	return &NodeStatus{
		Code:      statusMap[str],
		Msg:       str,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Status:    t.Status,
	}
}
