package main

import (
	"fmt"
  "net"
	"log"
	"net/url"
  "strconv"
  "strings"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
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

func mqttRpcRecv(client mqtt.Client, topic string, qos byte, ch chan string) {
  logger := log.New(os.Stdout, "[Mqtt recv] ", log.LstdFlags)
	client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
    m := string(msg.Payload())
    logger.Println(m)
    ch <- m
	})
}

func mqttRpcSend(client mqtt.Client, topic string, qos byte, ch chan string) {
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

func msgCenter(server Server) {

  Center := log.New(os.Stdout, "[Center] ", log.LstdFlags)
  Default := log.New(os.Stdout, "[Default] ", log.LstdFlags)

  input := make(chan string)

  // Mqtt
	//uri, err := url.Parse(os.Getenv("MQTT_URL"))
  uri, err := url.Parse(server.Mqtt)
  if err != nil {
    log.Fatal(err)
  }

  client := connect("node-" + strconv.Itoa(server.Id), uri)
	go mqttRpcRecv(client, "nodes/" + strconv.Itoa(server.Id) + "/rpc/send", 2, input)
  ch_mqtt := make(chan string, 100)
	go mqttRpcSend(client, "nodes/" + strconv.Itoa(server.Id) + "/rpc/recv", 2, ch_mqtt)

  // Ncp
  ch_ncp := make(chan string, 100)
  go ncp(ch_ncp, input)

  // Socket Client
  ch_socketc := make(chan string, 100)
  go socketClient(server.Tcpc, ch_socketc, input)

  // Router
  for {
    x := <- input
    Center.Println(x)

    switch {
    case isNcp(x):
      ch_ncp <- x
    case isJSONRPCRecv(x):
      ch_mqtt <- x
    case isJSONRPCSend(x):
      ch_socketc <- x
    default:
      Default.Println(x)
    }
  }
}

func socketClient(host string, input chan string, output chan string) {
  for {
    conn, err := net.Dial("tcp", host)
    //defer conn.Close()
    if err != nil {
      log.Println(err)
      time.Sleep(1000000000)
      // handle error
    } else {
      go socketSend(conn, input)
      socketRecv(conn, output)
    }
  }
}

func socketRecv(conn net.Conn, ch chan string) {
  logger := log.New(os.Stdout, "[Socket recv] ", log.LstdFlags)
  buf := make([]byte, 4096)
  for {
    cnt, err := conn.Read(buf)
    if err != nil || cnt == 0 {
      fmt.Println("EEEEEEEEEEee")
      break
    }
    inStr := strings.TrimSpace(string(buf[0:cnt]))
    logger.Println(inStr)
    ch <- inStr
  }
}

func socketSend(conn net.Conn, ch chan string) {
  logger := log.New(os.Stdout, "[Socket send] ", log.LstdFlags)
  for {
    m := <- ch
    logger.Println(m)
    _, e := conn.Write([]byte(m + "\n"))

    if e != nil {
      log.Printf("Error: ")
      log.Println(e)
      ch <- m
      break
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

  go msgCenter(config.Server)

  for {
    time.Sleep(1000000000)
  }

}
