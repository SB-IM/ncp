package main

import (
  "fmt"
  "log"
  "net/url"
  "strconv"
  "encoding/json"
  "os"
  "os/signal"
  "time"

  mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttRecv(client mqtt.Client, topic string, qos byte, ch chan string) {
  logger := log.New(os.Stdout, "[Mqtt recv] ", log.LstdFlags)
  rpc_filter := DuplicateFilter{}
  client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
    //fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
    m := string(msg.Payload())
    logger.Println(m)
    if rpc_filter.Put(m) != "" {
      ch <- m
    }
  })
}

func mqttSend(client mqtt.Client, topic string, qos byte, ch chan string) {
  logger := log.New(os.Stdout, "[Mqtt Send] ", log.LstdFlags)
  for x := range ch {
    logger.Println(x)
    client.Publish(topic, qos, true, x)
  }
}

type OnStatus struct {
  Code int `json:"code"`
  Msg string `json:"msg"`
  Timestamp string `json:"timestamp"`
  Status Status `json:"status"`
}

func mqttSetOnline(client mqtt.Client, status Status, topic string, s string) {
  statusMap := map[string]int {
    "online": 0,
    "offline": 1,
    "neterror": 2,
  }

  onstatus := &OnStatus{
    Code : statusMap[s],
    Msg : s,
    Timestamp : strconv.FormatInt(time.Now().Unix(), 10),
    Status : status,
  }

  r, _ := json.Marshal(onstatus)
  client.Publish(topic, 2, true, string(r))
}

func ncp(input chan string, output chan string) {
  logger := log.New(os.Stdout, "[Ncp cmd] ", log.LstdFlags)
  for cmd := range input {
    logger.Println(cmd)
    //output <- cmd
  }
}

func msgCenter(s chan os.Signal, server Server, n Ncp) {

  Center := log.New(os.Stdout, "[Center] ", log.LstdFlags)
  Default := log.New(os.Stdout, "[Default] ", log.LstdFlags)

  input := make(chan string)

  // Mqtt
  //uri, err := url.Parse(os.Getenv("MQTT_URL"))
  uri, err := url.Parse(server.Mqtt)
  if err != nil {
    log.Fatal(err)
  }

  //logger_mqtt := log.New(os.Stdout, "[Mqtt] ", log.LstdFlags)

  ch_mqtt := make(chan string, 100)

  mqtt := mqttProxy{
    id: strconv.Itoa(server.Id),
    ch_rpc_send: ch_mqtt,
    ch_rpc_recv: input,
  }

  //client := connect("node-" + strconv.Itoa(server.Id), uri, "nodes/" + strconv.Itoa(server.Id) + "/status", input)
  mqtt.Connect(n.Status, "node-" + strconv.Itoa(server.Id), uri, "nodes/" + strconv.Itoa(server.Id) + "/status")
  //go mqttRecv(client, "nodes/" + strconv.Itoa(server.Id) + "/rpc/send", 2, input)
  go mqttSend(mqtt.client, "nodes/" + strconv.Itoa(server.Id) + "/rpc/recv", 2, ch_mqtt)
  ch_mqtt_message := make(chan string, 100)
  go mqttSend(mqtt.client, "nodes/" + strconv.Itoa(server.Id) + "/message", 0, ch_mqtt_message)

  go func(){
    for sig := range s {
      // sig is a ^C, handle it
      fmt.Println("Got signal:", sig)
      mqttSetOnline(mqtt.client, n.Status, "nodes/" + strconv.Itoa(server.Id) + "/status", "offline")
      fmt.Println("set offline done")
      time.Sleep(10 * time.Millisecond)
      mqtt.client.Disconnect(1)
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

  // Socket tran
  go socketServerTran(server.Tran, mqtt.client, "nodes/" + strconv.Itoa(server.Id))

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
  config_path := "./config.yml"
  if os.Getenv("NCP_CONF") != "" {
    config_path = os.Getenv("NCP_CONF")
  }
  fmt.Println("load config: " + config_path)

  config, err := getConfig(config_path)
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
  go msgCenter(s, config.Server, config.Ncp)

  //for {
  //  time.Sleep(1000000000)
  //}
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  // Block until a signal is received.
  s <- <-c
  fmt.Println("Got signal:", s)
  time.Sleep(100 * time.Millisecond)
}

