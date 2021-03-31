package network

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"sb.im/ncp/ncpio"
	"sb.im/ncp/tests/help"
)

func cmdRun(str string) {
	fmt.Println("EXEC:", str)
	if out, err := help.CmdRun(str); err != nil {
		fmt.Printf("%s", out)
		panic(err)
	}
}

// Disable This test
// This test restore session
// Need:
// - mqtt sessionExpiryInterval
// - CleanStart:  false
func testBroker(t *testing.T) {
//func TestBroker(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
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
	help.GenerateMqttConfig(mqttdConfigPath, []byte(mqttdConfig))

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

	// Network simulation
	// tc qdisc add dev eth0 root netem delay 100ms 1000ms reorder 25% 50%
	// tc qdisc change dev eth0 root netem loss 20%
	// tc qdisc del dev eth0 root

	// Initial
	cmdRun("tc qdisc add dev eth0 root netem loss 0%")
	defer cmdRun("tc qdisc del dev eth0 root")

	cmdRun("tc qdisc change dev eth0 root netem loss 100%")

	ncpio.I <- []byte(fmt.Sprintf(`{"jsonrpc":"2.0","result":"ok","id":"test-m.%d"}`, 1))
	time.Sleep(10 * time.Second)
	ncpio.I <- []byte(fmt.Sprintf(`{"jsonrpc":"2.0","result":"ok","id":"test-m.%d"}`, 2))
	time.Sleep(10 * time.Second)

	cmdRun("tc qdisc change dev eth0 root netem loss 0%")

	ncpio.I <- []byte(fmt.Sprintf(`{"jsonrpc":"2.0","result":"ok","id":"test-m.%d"}`, 6))
	time.Sleep(3 * time.Second)
	fmt.Println("send msg")

	count := 1
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
}
