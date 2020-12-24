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
	id := "000"
	mqttRpcRecv, mqttRpcSend := "nodes/%s/rpc/recv", "nodes/%s/rpc/send"
	mqttAddr := "mqtt://localhost:1883"
	if addr := os.Getenv("MQTT"); addr != "" {
		// addr "localhost:1883"
		mqttAddr = fmt.Sprintf("mqtt://%s", addr)
	}
	var mqttdConfig = `
mqttd:
  id: ` + id + `
  static:
    link_id: 1
    lat: "22.6876423001"
    lng: "114.2248673001"
    alt: "10088.0001"
  client: "node-%s"
  status:  "nodes/%s/status"
  network: "nodes/%s/network"
  broker: ` + mqttAddr + `
  rpc :
    i: ` + mqttRpcRecv + `
    o: ` + mqttRpcSend + `
  gtran:
    prefix: "nodes/%s/msg/%s"
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

		pub := exec.CommandContext(ctx, "mosquitto_pub", "-L", mqttAddr+"/"+fmt.Sprintf(mqttRpcSend, id), "-m", req[i])
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
	sub := exec.CommandContext(ctx, "mosquitto_sub", "-L", mqttAddr+"/"+fmt.Sprintf(mqttRpcRecv, id))
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
