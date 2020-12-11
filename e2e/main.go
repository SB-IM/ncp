package main

import (
	"bufio"
	"context"
	"io/ioutil"
	"log"
	"net"
	"strings"

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

func main() {
	config, err := getConfig("e2e/config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("%+v\n", config)

	ncpios := ncpio.NewNcpIOs(config.NcpIO)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ncpios.Run(ctx)

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

	log.Println("Successfully")
}
