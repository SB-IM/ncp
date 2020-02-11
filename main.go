package main

import (
	"flag"
  "fmt"
  "log"
  "net/url"
  "strconv"
  "encoding/json"
  "os"
  "os/signal"
  "time"
	"regexp"
  //"strings"
  "reflect"

  mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

func ncp(ncp *NcpCmd, input chan string, output chan string) {
  logger := log.New(os.Stdout, "[Ncp cmd] ", log.LstdFlags)
  for cmd := range input {
    logger.Println(cmd)
    output <- ncpCmd(ncp, cmd)
    //output <- cmd
  }
}

func msgCenter(s chan os.Signal, server Server, ncpCmd *NcpCmd, n Ncp, config_log *ConfigLog) {
	logGroup, err := logGroupNew(config_log)
	if err != nil {
		log.Fatal(err)
	}

	// default input
	input := make(chan string, 100)

  // Mqtt
  //uri, err := url.Parse(os.Getenv("MQTT_URL"))
  uri, err := url.Parse(server.Mqtt)
  if err != nil {
    log.Fatal(err)
  }

  ch_mqtt_i := make(chan string, 100)
  ch_mqtt_o := make(chan string, 100)

  mqtt := mqttProxy{
    id: strconv.Itoa(server.Id),
    ch_rpc_send: ch_mqtt_o,
    ch_rpc_recv: ch_mqtt_i,
  }


	mqtt.Connect(n.Status, logGroup.Get("mqtt"), "node-" + strconv.Itoa(server.Id), uri, "nodes/" + strconv.Itoa(server.Id) + "/status")
	go mqttSend(mqtt.client, logGroup.Get("mqtt"), "nodes/" + strconv.Itoa(server.Id) + "/rpc/recv", 2, true, ch_mqtt_o)

	ch_mqtt_msg := make(chan string, 100)
	go mqttTran(mqtt.client, logGroup.Get("mqtr"), "nodes/" + strconv.Itoa(server.Id), ch_mqtt_msg)

	defer func() {
		mqttSetOnline(mqtt.client, n.Status, "nodes/" + strconv.Itoa(server.Id) + "/status", "offline")
		fmt.Println("set offline done")

		// Need Wait Set offline
		time.Sleep(10 * time.Millisecond)
		mqtt.client.Disconnect(1)

		if err = logGroup.Close(); err != nil {
			log.Fatal(err)
		}
	}()

  // Ncp
  ch_ncp := make(chan string, 100)
  go ncp(ncpCmd, ch_ncp, input)

  // Socket Client
  ch_socketc_i := make(chan string, 100)
  ch_socketc := make(chan string, 100)
	socketClient := &SocketClient{
		logger: logGroup.Get("tcpc"),
	}
	go socketClient.Run(server.Tcpc, ch_socketc, ch_socketc_i)

	// Socket Server
	ch_sockets := make(chan string, 100)
	socketServer := &SocketServer{
		logger: logGroup.Get("tcps"),
	}
	go socketServer.Listen(server.Tcps, ch_sockets, input)

	// Socket tran
	ch_no_message := make(chan string)
	socketServerTran := &SocketServer{
		logger: logGroup.Get("tran"),
	}
	go socketServerTran.Listen(server.Tran, ch_no_message, input)

	// Router
	var x string
	rpc_run := RpcRun{}
	Filter := log.New(os.Stdout, "[Filter] ", log.LstdFlags)

  for {
    select {
    case x = <- ch_mqtt_i:
			if rpc_run.Run(x) {
				// Disable Notify ack
				//ch_mqtt_o <- confirmNotice(x)
			} else {
				Filter.Println(x)
				x = ""
			}
		case x = <-ch_socketc_i:
			if isLink(x) {
				bit13_timestamp := string([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))[:13])
				override_id := "ncp." + strconv.Itoa(server.Id) + "-" + bit13_timestamp

				req, err, callback := linkCall([]byte(x), override_id)
				if err != nil {
					fmt.Println(err)
				}

				go func() {
					res, err := callback([]byte(syncMqttRpc(mqtt.client, logGroup.Get("mqtt"), server.Link_id, string(req))))
					if err != nil {
						fmt.Println(err)
					}
					ch_socketc <-string(res)
				}()

				fmt.Println("[Link Call] " + x)
				x = ""
			}

    case x = <- input:
      //fmt.Println("Recvice", x)
		case sig := <- s:
			// sig is a ^C, handle it
			fmt.Println("Got signal:", sig)
			return
    }

    if x == "" { continue }

    //Center.Println(x)

    switch {
    case isNcp(x):
      ch_ncp <- x
      ch_sockets <- x
    case isJSONRPCRecv(x):
      ch_mqtt_o <- x
    case isJSONRPCSend(x):
      ch_socketc <- x
      ch_sockets <- x
    default:
      //ch_mqtt_message <- x
      ch_mqtt_msg <- x
    }
  }
}

func ncpCmd(ncp *NcpCmd, raw string) string {
	rpc := getJSONRPC(raw)
  results := CallObjectMethod(ncp, Ucfirst("status"))

	if regexp.MustCompile(`^\{.*\}$`).MatchString(string(*rpc.Params)) {
		results = CallObjectMethod(ncp, Ucfirst("webrtc"), *rpc.Params)

	} else if rpc.Method == "ncp" {
		fmt.Println(string(*rpc.Params))

		var params []string
		json.Unmarshal(*rpc.Params, &params)

		switch params[0] {
		case "status":
			results = CallObjectMethod(ncp, Ucfirst("status"))
		case "upload":
			results = CallObjectMethod(ncp, Ucfirst("upload"), params[1], params[2])
		case "download":
			//results := CallObjectMethod(ncp, Ucfirst("download"), "map", "http://localhost:3000/ncp/v1/plans/12/get_map")
			results = CallObjectMethod(ncp, Ucfirst("download"), params[1], params[2])
		case "shell":
			results = CallObjectMethod(ncp, Ucfirst("shell"), params[1])
		default:
			results = CallObjectMethod(ncp, Ucfirst("status"))
		}
	} else {
		var params []string
		json.Unmarshal(*rpc.Params, &params)
		switch rpc.Method {
		case "status":
			results = CallObjectMethod(ncp, Ucfirst("status"))
		case "upload":
			results = CallObjectMethod(ncp, Ucfirst("upload"), params[0], params[1])
		case "download":
			results = CallObjectMethod(ncp, Ucfirst("download"), params[0], params[1])
		case "shell":
			results = CallObjectMethod(ncp, Ucfirst("shell"), params[0])
		default:
			results = CallObjectMethod(ncp, Ucfirst("status"))
		}
	}


  //results := CallObjectMethod(ncp, Ucfirst(rpc.Method))

	//tmp := CallObjectMethod(new(Ncp), "Method_" + method)
  //fmt.Printf(strings.TrimSuffix(strings.TrimPrefix(str, "["), "]"))

	result := results.([]reflect.Value)[0].Interface().([]byte)

	var s string
	if string(result) == "" {
		s =`"result":""`
	} else {
		s =`"result":` + string(result)
	}


	fmt.Println(string(result))
	if e := results.([]reflect.Value)[1].Interface(); e != nil {
		fmt.Println(e.(error))
		s = `"error": "` + e.(error).Error() + `"`
	}

	return `{"jsonrpc":"2.0",`+s+`,"id":"` + rpc.Id + `"}`
}

func main() {
	config_path := "config.yml"

	help := flag.Bool("h", false, "this help")
	flag.StringVar(&config_path, "c", "config.yml", "set configuration file")

	show_version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *show_version {
		fmt.Println(version)
		return
	}

  if os.Getenv("NCP_CONF") != "" {
    config_path = os.Getenv("NCP_CONF")
  }
  fmt.Println("load config: " + config_path)

  config, err := getConfig(config_path)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  fmt.Println("=========")
  //fmt.Println(config)

  config.Ncp.Common.Id = strconv.Itoa(config.Server.Id)
  config.Ncp.Common.SecretKey = config.Server.Secret_key

  config.Server.Link_id = config.Ncp.Status.Link_id
  ncpCmd := NcpCmd {
    config: config.Ncp,
  }

	// init, Golang has no constructor
  ncpCmd.Init()

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

	msgCenter(c, config.Server, &ncpCmd, config.Ncp, &config.Log)
}

