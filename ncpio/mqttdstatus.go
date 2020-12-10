package ncpio

import (
	"time"

	"github.com/go-ping/ping"

	logger "log"
)

type StaticStatus struct {
	Link_id     int    `json:"link_id"`
	Position_ok bool   `json:"position_ok"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	Alt         string `json:"alt"`
}

type NetworkStatus struct {
	// (0 ~ 100)%
	Loss int `json:"loss"`
	// AvgRtt (ms)
	Time int `json:"time"`
}

// {"code":0,"msg":"online","timestamp":"1607566337","status":{"link_id":3,"position_ok":true,"lat":"22.831878298","lng":"113.514688221","alt":"80.0001"}}
type NodeStatus struct {
	Code      int           `json:"code"`
	Msg       string        `json:"msg"`
	Timestamp string        `json:"timestamp"`
	Status    *StaticStatus `json:"status"`
}

func NetworkPing(addr string, callback func(*NetworkStatus)) {
	for {
		pinger, err := ping.NewPinger(addr)
		if err != nil {
			logger.Println(err)
		}
		pinger.OnFinish = func(stats *ping.Statistics) {
			callback(&NetworkStatus{
				Loss: int(stats.PacketLoss),
				Time: int(stats.AvgRtt.Milliseconds()),
			})
		}

		pinger.Count = 3
		err = pinger.Run()
		if err != nil {
			logger.Println(err)
		}
		time.Sleep(60 * time.Second)
	}
}
