package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
	"sb.im/ncp/ncpio"
)

type Config struct {
	NcpIO []ncpio.Config `json:"ncpio"`
}

func getConfig(str string) (Config, error) {
	config := Config{}
	configFile, err := ioutil.ReadFile(str)
	if err != nil {
		return config, err
	} else {
		err = yaml.Unmarshal(configFile, &config)
		return config, err
	}
}

func generateMqttConfig(name string, config []byte) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	if _, err := file.Write(config); err != nil {
		panic(err)
	}
}

func broker() {
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
	tmpConfig := "test_mqttd.yml"
	generateMqttConfig(tmpConfig, []byte(mqttdConfig))

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config, err := getConfig("e2e/e2e.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("%+v\n", config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ncpios := ncpio.NewNcpIOs(config.NcpIO)
	go ncpios.Run(ctx)

	// wait tcp server startup
	time.Sleep(3 * time.Millisecond)

	msg1 := "2333333333333"
	msg2 := "4555555555555"

	ncpio.I <- []byte(msg1)
	ncpio.I <- []byte(msg2)

	addr := "localhost:1222"
	for _, c := range config.NcpIO {
		if c.Type == "tcps" {
			addr = c.Params
		}
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	if strings.TrimSuffix(status, "\n") != msg1 {
		log.Panicln("Should", msg1)
	}

	status, err = reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	if strings.TrimSuffix(status, "\n") != msg2 {
		log.Panicln("Should", msg2)
	}

	conn.Write([]byte(msg2 + "\n"))
	conn.Write([]byte(msg1 + "\n"))

	if string(<-ncpio.O) != msg2 {
		log.Panicln("Should", msg2)
	}

	if string(<-ncpio.O) != msg1 {
		log.Panicln("Should", msg1)
	}

	log.Println("Socket Successfully")

	topic := "nodes/999/rpc/send"
	pub := exec.CommandContext(ctx, "mosquitto_pub", "-L", mqttAddr+"/"+topic, "-m", "xxxxx")
	if data, err := pub.CombinedOutput(); err != nil {
		log.Printf("%s\n", data)
		log.Panicln(err)
	}

	log.Printf("%s\n", <-ncpio.O)

	os.Remove(tmpConfig)
	log.Println("Successfully")
}
