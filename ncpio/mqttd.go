package ncpio

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	//"os"
	"time"

	"sb.im/ncp/util"

	"github.com/SB-IM/jsonrpc-lite"
	packets "github.com/eclipse/paho.golang/packets"
	paho "github.com/eclipse/paho.golang/paho"

	logger "log"
)

type Mqtt struct {
	Client  *paho.Client
	Connect *paho.Connect
	Config  *MqttdConfig
	status  *NodeStatus
	cache   chan []byte
	I       <-chan []byte
	O       chan<- []byte
}

func NewMqtt(params string, i <-chan []byte, o chan<- []byte) *Mqtt {
	config, err := loadMqttConfigFromFile(params)
	if err != nil {
		logger.Println(err)
	}

	opt, err := url.Parse(config.Broker)
	if err != nil {
		logger.Println(err)
		return nil
	}
	logger.Printf("%+v\n", config)
	logger.Printf("%+v\n", config.Static)

	password, _ := opt.User.Password()
	status := &NodeStatus{
		Status: config.Static,
	}
	raw, _ := json.Marshal(status.SetOnline("neterror"))
	return &Mqtt{
		I:      i,
		O:      o,
		cache:  make(chan []byte, 128),
		status: status,
		Config: config,
		Client: paho.NewClient(paho.ClientConfig{
			ClientID: fmt.Sprint(config.Client, config.ID),
			Router: paho.NewSingleHandlerRouter(func(p *paho.Publish) {
				o <- p.Payload
			}),
		}),
		Connect: paho.ConnectFromPacketConnect(&packets.Connect{
			WillProperties: &packets.Properties{},

			WillFlag:    true,
			WillMessage: raw,
			WillRetain:  true,
			WillTopic:   fmt.Sprintf(config.Status, config.ID),
			WillQOS:     1,
			Password:    []byte(password),
			Username:    opt.User.Username(),
			ClientID:    fmt.Sprintf(config.Client, config.ID),
			CleanStart:  false,
			// interval 10s
			KeepAlive: 10,
			// TODO:
			Properties: &packets.Properties{
				// PayloadFormat indicates the format of the payload of the message
				// 0 is unspecified bytes
				// 1 is UTF8 encoded character data
				//PayloadFormat: 1,
				// MessageExpiry is the lifetime of the message in seconds
				//MessageExpiry *uint32
				//// ContentType is a UTF8 string describing the content of the message
				//// for example it could be a MIME type
				//ContentType string
				//// ResponseTopic is a UTF8 string indicating the topic name to which any
				//// response to this message should be sent
				//ResponseTopic string
				//// CorrelationData is binary data used to associate future response
				//// messages with the original request message
				//CorrelationData []byte
				//// SubscriptionIdentifier is an identifier of the subscription to which
				//// the Publish matched
				//SubscriptionIdentifier *uint32
				//// SessionExpiryInterval is the time in seconds after a client disconnects
				//// that the server should retain the session information (subscriptions etc)
				//SessionExpiryInterval *uint32
				//// AssignedClientID is the server assigned client identifier in the case
				//// that a client connected without specifying a clientID the server
				//// generates one and returns it in the Connack
				//AssignedClientID string
				//// ServerKeepAlive allows the server to specify in the Connack packet
				//// the time in seconds to be used as the keep alive value
				//ServerKeepAlive *uint16
				//// AuthMethod is a UTF8 string containing the name of the authentication
				//// method to be used for extended authentication
				//AuthMethod string
				//// AuthData is binary data containing authentication data
				//AuthData []byte
				//// RequestProblemInfo is used by the Client to indicate to the server to
				//// include the Reason String and/or User Properties in case of failures
				//RequestProblemInfo *byte
				//// WillDelayInterval is the number of seconds the server waits after the
				//// point at which it would otherwise send the will message before sending
				//// it. The client reconnecting before that time expires causes the server
				//// to cancel sending the will
				//WillDelayInterval *uint32
				//// RequestResponseInfo is used by the Client to request the Server provide
				//// Response Information in the Connack
				//RequestResponseInfo *byte
				//// ResponseInfo is a UTF8 encoded string that can be used as the basis for
				//// createing a Response Topic. The way in which the Client creates a
				//// Response Topic from the Response Information is not defined. A common
				//// use of this is to pass a globally unique portion of the topic tree which
				//// is reserved for this Client for at least the lifetime of its Session. This
				//// often cannot just be a random name as both the requesting Client and the
				//// responding Client need to be authorized to use it. It is normal to use this
				//// as the root of a topic tree for a particular Client. For the Server to
				//// return this information, it normally needs to be correctly configured.
				//// Using this mechanism allows this configuration to be done once in the
				//// Server rather than in each Client
				//ResponseInfo string
				//// ServerReference is a UTF8 string indicating another server the client
				//// can use
				//ServerReference string
				//// ReasonString is a UTF8 string representing the reason associated with
				//// this response, intended to be human readable for diagnostic purposes
				//ReasonString string
				//// ReceiveMaximum is the maximum number of QOS1 & 2 messages allowed to be
				//// 'inflight' (not having received a PUBACK/PUBCOMP response for)
				//ReceiveMaximum *uint16
				//// TopicAliasMaximum is the highest value permitted as a Topic Alias
				//TopicAliasMaximum *uint16
				//// TopicAlias is used in place of the topic string to reduce the size of
				//// packets for repeated messages on a topic
				//TopicAlias *uint16
				//// MaximumQOS is the highest QOS level permitted for a Publish
				//MaximumQOS *byte
				//// RetainAvailable indicates whether the server supports messages with the
				//// retain flag set
				//RetainAvailable *byte
				//// User is a map of user provided properties
				//User map[string]string
				//// MaximumPacketSize allows the client or server to specify the maximum packet
				//// size in bytes that they support
				//MaximumPacketSize *uint32
				//// WildcardSubAvailable indicates whether wildcard subscriptions are permitted
				//WildcardSubAvailable *byte
				//// SubIDAvailable indicates whether subscription identifiers are supported
				//SubIDAvailable *byte
				//// SharedSubAvailable indicates whether shared subscriptions are supported
				//SharedSubAvailable *byte
			},
		}),
	}
}

func (t *Mqtt) Run(ctx context.Context) {
	opt, err := url.Parse(t.Config.Broker)
	if err != nil {
		logger.Println(err)
	}

	for {
		logger.Println("MQTT Connect")
		if conn, err := net.Dial("tcp", opt.Hostname()+":"+opt.Port()); err != nil {
			logger.Println(err)
		} else {
			t.Client.Conn = conn
			t.doRun(ctx)
			conn.Close()
		}
		time.Sleep(1 * time.Second)
	}
}

func (t *Mqtt) doRun(ctx context.Context) {
	pinger := NewPingHandler(t.Client, fmt.Sprintf(t.Config.Network, t.Config.ID))
	//pinger.SetDebug(logger.New(os.Stdout, "[Pinger]: ", logger.LstdFlags | logger.Lshortfile))
	t.Client.PingHandler = pinger

	defer logger.Println("MQTT exit")
	t.Client.Connect(ctx, t.Connect)

	res, err := t.Client.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			fmt.Sprintf(t.Config.Rpc.O, t.Config.ID): paho.SubscribeOptions{
				QoS: 2,
				//RetainHandling    byte
				//NoLocal           bool
				//RetainAsPublished bool
			},
		},
	})

	if err != nil {
		if res != nil {
			logger.Printf("%+v\n", res)
		}
		logger.Println(err)
		return
	}

	defer t.setStatus("offline")
	t.setStatus("online")

	for {
		select {
		case raw := <-t.cache:
			if err := t.send(ctx, raw); err != nil {
				t.cache <- raw
				return
			}
		case raw := <-t.I:
			if err := t.send(ctx, raw); err != nil {
				t.cache <- raw
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (t *Mqtt) send(ctx context.Context, raw []byte) error {
	if rpc, err := jsonrpc.Parse(raw); err == nil && (rpc.Type == jsonrpc.TypeSuccess || rpc.Type == jsonrpc.TypeErrors) {
		//fmt.Println("[RES]: ", string(raw))
		// {"jsonrpc":"2.0","result":"ok","id":"test.0-1607482556696-0"}
		// {"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":"test.0-99991607483766.0"}

		res, err := t.Client.Publish(ctx, &paho.Publish{
			Payload: raw,
			Topic:   fmt.Sprintf(t.Config.Rpc.I, t.Config.ID),
			QoS:     2,
		})

		if err != nil {
			if res != nil {
				logger.Printf("%+v\n", res)
			}
			logger.Println(err)
			return err
		}
	} else if err == nil && (rpc.Type == jsonrpc.TypeRequest || rpc.Type == jsonrpc.TypeNotify) {
		//fmt.Println("[REQ]: ", string(raw))
		// JSON-RPC Request Ignore

		// {"jsonrpc":"2.0","method":"test","params":[]}
		// {"jsonrpc":"2.0","id":"test.0-1553321035000","method":"test","params":[]}

		// {"jsonrpc":"2.0","method":"ncp_offline"}
		onlineStatus := "offline"
		if rpc.Method == "ncp_online" {
			onlineStatus = "online"
		}

		t.setStatus(onlineStatus)
	} else {
		//fmt.Println("[Tran]: ", string(raw))

		for key, data := range util.DetachTran(raw) {
			opt, ok := t.Config.Trans[key]
			if !ok {
				// TODO: 'opt' use Default
			}
			res, err := t.Client.Publish(ctx, &paho.Publish{
				Payload: data,
				Topic:   fmt.Sprintf(t.Config.Gtran.Prefix, t.Config.ID, key),
				QoS:     opt.QoS,
				Retain:  opt.Retain,
				//Properties *Properties
				//PacketID   uint16
				//Duplicate  bool

			})

			if err != nil {
				if res != nil {
					logger.Printf("%+v\n", res)
				}
				logger.Println(err)
				return err
			}
		}
	}
	return nil
}

func (t *Mqtt) setStatus(str string) error {
	raw, err := json.Marshal(t.status.SetOnline(str))
	if err != nil {
		logger.Println(err)
		return err
	} else {
		res, err := t.Client.Publish(context.Background(), &paho.Publish{
			Payload: raw,
			Topic:   fmt.Sprintf(t.Config.Status, t.Config.ID),
			QoS:     1,
			Retain:  true,
		})

		if err != nil {
			if res != nil {
				logger.Printf("%+v\n", res)
			}
			logger.Println(err)
			return err
		}
	}
	return nil
}
