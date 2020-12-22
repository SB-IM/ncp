package ncpio

import (
	"context"
	"net"
	"time"

	logger "log"
)

func (t *Tcpc) Listen(ctx context.Context) {
	ln, err := net.Listen("tcp", t.addr)
	if err != nil {
		logger.Println("TCPS:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			time.Sleep(retryInterval)
		} else {
			logger.Println("New Tcps Client Connect")

			t.ctx, t.cancel = context.WithCancel(context.Background())
			go t.send(conn)
			go t.recv(conn)
			select {
			case <-ctx.Done():
				if err := ln.Close(); err != nil {
					logger.Println("Connect Close error:", err)
				}
				return
			case <-t.ctx.Done():
				conn.Close()
				logger.Println("Connect err try reconnect")
			}
		}
	}
}
