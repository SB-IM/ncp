package ncpio

// Fork from: https://github.com/eclipse/paho.golang/blob/master/paho/pinger.go

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	packets "github.com/eclipse/paho.golang/packets"
	paho "github.com/eclipse/paho.golang/paho"
)

type PingHandler struct {
	client  *paho.Client
	topic   string
	lastPub string
	time    time.Time
	delay   time.Duration

	mu              sync.Mutex
	lastPing        time.Time
	conn            net.Conn
	stop            chan struct{}
	pingFailHandler paho.PingFailHandler
	pingOutstanding int32
	debug           paho.Logger
}

func NewPingHandler(client *paho.Client, topic string) *PingHandler {
	return &PingHandler{
		client:          client,
		topic:           topic,
		pingFailHandler: func(e error) {},
		debug:           paho.NOOPLogger{},
	}
}

// Start is the library provided Pinger's implementation of
// the required interface function()
func (p *PingHandler) Start(c net.Conn, pt time.Duration) {
	p.mu.Lock()
	p.conn = c
	p.stop = make(chan struct{})
	p.mu.Unlock()
	checkTicker := time.NewTicker(pt / 4)
	defer checkTicker.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-checkTicker.C:
			if atomic.LoadInt32(&p.pingOutstanding) > 0 && time.Since(p.lastPing) > (pt+pt>>1) {
				p.pingFailHandler(fmt.Errorf("ping resp timed out"))
				//ping outstanding and not reset in 1.5 times ping timer
				return
			}
			if time.Since(p.lastPing) >= pt {
				//time to send a ping
				if _, err := packets.NewControlPacket(packets.PINGREQ).WriteTo(p.conn); err != nil {

					p.debug.Println("pingHandler sending ping request")
					if p.pingFailHandler != nil {
						p.pingFailHandler(err)
					}
					return
				}
			}
			atomic.AddInt32(&p.pingOutstanding, 1)
			p.time = time.Now()
			p.debug.Println("pingHandler sending ping request")

		}
	}
}

// Stop is the library provided Pinger's implementation of
// the required interface function()
func (p *PingHandler) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stop == nil {
		return
	}
	p.debug.Println("pingHandler stopping")
	select {
	case <-p.stop:
		//Already stopped, do nothing
	default:
		close(p.stop)
	}
}

// PingResp is the library provided Pinger's implementation of
// the required interface function()
func (p *PingHandler) PingResp() {
	delay := fmt.Sprintf("%d", time.Since(p.time).Milliseconds())
	p.debug.Printf("delay: %s, pingHandler resetting pingOutstanding", time.Since(p.time))
	atomic.StoreInt32(&p.pingOutstanding, 0)

	if p.lastPub != delay {
		res, err := p.client.Publish(context.TODO(), &paho.Publish{
			Payload: []byte(delay),
			Topic:   p.topic,
			QoS:     0,
			Retain:  true,
		})

		if err != nil {
			if res != nil {
				p.debug.Printf("%+v\n", res)
			}
			p.debug.Println(err)
			return
		}
	}
	p.lastPub = delay
}

// SetDebug sets the logger l to be used for printing debug
// information for the pinger
func (p *PingHandler) SetDebug(l paho.Logger) {
	p.debug = l
}
