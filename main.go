package main

import (
	"fmt"
	"log"
	"net/url"
  "strconv"
  "os"
  "os/signal"
  "time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientId string, uri *url.URL, willTopic string) mqtt.Client {
	opts := createClientOptions(clientId, uri)
  opts.SetWill(willTopic, "2", 2, true)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func mqttRecv(client mqtt.Client, topic string, qos byte, ch chan string) {
  logger := log.New(os.Stdout, "[Mqtt recv] ", log.LstdFlags)
	client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
    m := string(msg.Payload())
    logger.Println(m)
    ch <- m
	})
}

func mqttSend(client mqtt.Client, topic string, qos byte, ch chan string) {
  logger := log.New(os.Stdout, "[Mqtt Send] ", log.LstdFlags)
  for x := range ch {
    logger.Println(x)
    client.Publish(topic, qos, true, x)
  }
}

func mqttSetOnline(client mqtt.Client, topic string, status string) {
  statusMap := map[string]string {
    "online": "0",
    "offline": "1",
    "neterror": "2",
  }

  client.Publish(topic, 2, true, statusMap[status])
}

func ncp(input chan string, output chan string) {
  logger := log.New(os.Stdout, "[Ncp cmd] ", log.LstdFlags)
  for cmd := range input {
    logger.Println(cmd)
    //output <- cmd
  }
}

func msgCenter(s chan os.Signal, server Server) {

  Center := log.New(os.Stdout, "[Center] ", log.LstdFlags)
  Default := log.New(os.Stdout, "[Default] ", log.LstdFlags)

  input := make(chan string)

  // Mqtt
	//uri, err := url.Parse(os.Getenv("MQTT_URL"))
  uri, err := url.Parse(server.Mqtt)
  if err != nil {
    log.Fatal(err)
  }

  client := connect("node-" + strconv.Itoa(server.Id), uri, "nodes/" + strconv.Itoa(server.Id) + "/status")
	go mqttRecv(client, "nodes/" + strconv.Itoa(server.Id) + "/rpc/send", 2, input)
  ch_mqtt := make(chan string, 100)
	go mqttSend(client, "nodes/" + strconv.Itoa(server.Id) + "/rpc/recv", 2, ch_mqtt)
  ch_mqtt_message := make(chan string, 100)
	go mqttSend(client, "nodes/" + strconv.Itoa(server.Id) + "/message", 0, ch_mqtt_message)

  mqttSetOnline(client, "nodes/" + strconv.Itoa(server.Id) + "/status", "online")

  go func(){
    for sig := range s {
      // sig is a ^C, handle it
      fmt.Println("Got signal:", sig)
      mqttSetOnline(client, "nodes/" + strconv.Itoa(server.Id) + "/status", "offline")
      fmt.Println("set offline done")
      time.Sleep(100000000)
      client.Disconnect(1)
    }
  }()

  // Ncp
  ch_ncp := make(chan string, 100)
  go ncp(ch_ncp, input)

  // Socket Client
  ch_socketc := make(chan string, 100)
  go socketClient(server.Tcpc, ch_socketc, input)

  // Socket Server
  ch_sockets := make(chan string, 100)
  go socketServer(server.Tcps, ch_sockets, input)

  // Router
  for {
    x := <- input
    Center.Println(x)

    switch {
    case isNcp(x):
      ch_ncp <- x
      ch_sockets <- x
    case isJSONRPCRecv(x):
      ch_mqtt <- x
    case isJSONRPCSend(x):
      ch_socketc <- x
    default:
      Default.Println(x)
      ch_mqtt_message <- x
    }
  }
}

func main() {
  config, err := getConfig("./config.yml")
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  fmt.Println("=========")
  //topic := "test"
	//topic := uri.Path[1:len(uri.Path)]
	//if topic == "" {
	//	topic = "test"
	//}

  s := make(chan os.Signal)
  go msgCenter(s, config.Server)

  //for {
  //  time.Sleep(1000000000)
  //}
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  // Block until a signal is received.
  s <- <-c
  fmt.Println("Got signal:", s)
  time.Sleep(1000000000)

}
