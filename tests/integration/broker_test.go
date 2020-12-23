package integration

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"sb.im/ncp/ncpio"
)

func TestBroker(t *testing.T) {
	mqttAddr := "mqtt://localhost:1883"
	if addr := os.Getenv("MQTT"); addr != "" {
		// addr "localhost:1883"
		mqttAddr = fmt.Sprintf("mqtt://%s", addr)
	}
	var mqttdConfig = `
mqttd:
  id: 999
  static:
    link_id: 1
    lat: "22.6876423001"
    lng: "114.2248673001"
    alt: "10088.0001"
  client: "node-%d"
  status:  "nodes/%d/status"
  network: "nodes/%d/network"
  broker: ` + mqttAddr + `
  rpc :
    i: "nodes/%d/rpc/recv"
    o: "nodes/%d/rpc/send"
  gtran:
    prefix: "nodes/%d/msg/%s"
  trans:
    wether:
      retain: true
      qos: 0
    battery:
      retain: true
      qos: 0

`
	mqttdConfigPath := "test_mqttd.yml"
	generateMqttConfig(mqttdConfigPath, []byte(mqttdConfig))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ncpios := ncpio.NewNcpIOs([]ncpio.Config{
		{
			Type: "api",
			IRules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
			ORules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
		},
		{
			Type:   "mqtt",
			Params: mqttdConfigPath,
			IRules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
			ORules: []ncpio.Rule{
				{
					Regexp: `.*`,
					Invert: false,
				},
			},
		},
	})
	go ncpios.Run(ctx)

	// Wait mqttd server startup && sub topic on broker
	time.Sleep(3 * time.Millisecond)

	// Wait load mqttdConfig after delete
	os.Remove(mqttdConfigPath)

	// mqttd buffer 128
	// api   buffer 128
	// Buffer Max 128 + 128 = 256
	//count := 260
	count := 0
	req := make([]string, count)

	// Send
	for i := 0; i < count; i++ {
		req[i] = fmt.Sprintf(`{"jsonrpc":"2.0","method":"test","id":"test.%d"}`, i)

		sendTopic := "nodes/999/rpc/send"
		pub := exec.CommandContext(ctx, "mosquitto_pub", "-L", mqttAddr+"/"+sendTopic, "-m", req[i])
		if out, err := pub.CombinedOutput(); err != nil {
			t.Error(string(out), err)
		}
	}

	for i := 0; i < count; i++ {
		select {
		case res := <-ncpio.O:
			if string(res) != req[i] {
				t.Errorf("Recv is: %s, Should: %s", res, req[i])
			}
		case <-time.After(1 * time.Millisecond):
			// more than buffer 256 is dropped
			if i <= 256 {
				t.Error("Recv timeout", i)
			}
		}
	}

	// Recv
	recvTopic := "nodes/999/rpc/recv"
	sub := exec.CommandContext(ctx, "mosquitto_sub", "-L", mqttAddr+"/"+recvTopic)
	stdout, err := sub.StdoutPipe()
	if err != nil {
		t.Error(err)
	}
	if err := sub.Start(); err != nil {
		t.Error(err)
	}
	// Wait sub topic on broker
	time.Sleep(10 * time.Millisecond)

	// api   buffer 128
	count = 128
	req = make([]string, count)
	for i := 0; i < count; i++ {
		req[i] = fmt.Sprintf(`{"jsonrpc":"2.0","result":"ok","id":"test.%d"}`, i)
		ncpio.I <- []byte(req[i])
	}

	reader := bufio.NewReader(stdout)
	for i := 0; i < count; i++ {
		raw, err := reader.ReadString('\n')
		if err != nil {
			t.Error(err)
		}
		if res := strings.TrimSuffix(raw, "\n"); res != req[i] {
			t.Errorf("Recv is: %s, Should: %s", res, req[i])
		}
	}
}
