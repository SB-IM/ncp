package main

import (
  "log"
  "time"
  "strconv"
  "net/url"

	"sb.im/ncp/history"

  mqtt "github.com/eclipse/paho.mqtt.golang"
)

func setUri(uri *url.URL) *mqtt.ClientOptions {
  opts := mqtt.NewClientOptions()
  //opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
  opts.AddBroker("tcp://" + uri.Host)
  opts.SetUsername(uri.User.Username())
  password, _ := uri.User.Password()
  opts.SetPassword(password)
  return opts
}

type mqttProxy struct {
  client mqtt.Client
  id string
  ch_rpc_send chan string
  ch_rpc_recv chan string
}

func (this *mqttProxy) Connect(status Status, logger *log.Logger, clientId string, uri *url.URL, willTopic string) mqtt.Client {
  opts := setUri(uri)
  opts.SetWill(willTopic, `{"code":2,"msg":"neterror"}`, 2, true)

  opts.SetClientID(clientId)

  // interval 2s
  opts.SetKeepAlive(2 * time.Second)
  opts.SetResumeSubs(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)

  // Lost callback
  opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
    logger.Println("Lost Connect")
  })

  // Connect && Reconnect callback
  opts.SetOnConnectHandler(func (client mqtt.Client) {
    clientOptionsReader := client.OptionsReader()
    mqttSetOnline(client, status, clientOptionsReader.WillTopic(), "online")

		logger.Println("New Connect", &client)
		go mqttRecv(client, logger, "nodes/" + this.id + "/rpc/send", 2, this.ch_rpc_recv)
  })

  client := mqtt.NewClient(opts)
  token := client.Connect()
  for !token.WaitTimeout(3 * time.Second) {
  }
  if err := token.Error(); err != nil {
    logger.Fatal(err)
  }
  this.client = client
  return client
}

// Refactor After ----------------

func syncMqttRpc(client mqtt.Client, logger *log.Logger, id int, send string) string {
	invoking := getJSONRPC(send).Id
	topic := "nodes/" + strconv.Itoa(id) + "/rpc/"
	ch_recv := make(chan string)
	client.Subscribe(topic + "recv", 1, func(client mqtt.Client, mqtt_msg mqtt.Message) {
		msg := string(mqtt_msg.Payload())
		if invoking == getJSONRPC(msg).Id {
			client.Unsubscribe(topic + "recv")
			logger.Println("Res:", msg)
			ch_recv <-msg
		}
	})

	msgSend := func() {
		client.Publish(topic + "send", 2, false, send)
		logger.Println("Req:", send)
	}

	msgSend()

	for {
		select {
		case <-time.After(10 * time.Second):
			msgSend()
		case result := <-ch_recv:
			return result
		}
	}
}

func mqttRecv(client mqtt.Client, logger *log.Logger, topic string, qos byte, ch chan string) {
	client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		m := string(msg.Payload())
		logger.Println("Recv:", m)
		ch <- m
	})
}

func mqttSend(client mqtt.Client, logger *log.Logger, topic string, qos byte, retained bool, ch chan string) {
	for x := range ch {
		logger.Println("Send:", x)
		client.Publish(topic, qos, retained, x)
	}
}

func mqttTran(archive *history.Archive, client mqtt.Client, logger *log.Logger, topic_prefix string, ch chan string) {
	for x := range ch {
		dataMap := detachTran([]byte(x))
		for k, v := range dataMap {
			if archive.FilterAdd(k, v) {
				logger.Println(k + " --> " + string(v))
				client.Publish(topic_prefix + "/msg/" + k, 0, true, string(v))
			}
		}
	}
}

