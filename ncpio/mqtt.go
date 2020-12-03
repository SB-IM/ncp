package ncpio

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"sb.im/ncp/util"

	packets "github.com/eclipse/paho.golang/packets"
	paho "github.com/eclipse/paho.golang/paho"
)

type Mqtt struct {
	Client  *paho.Client
	Connect *paho.Connect
	I       <-chan []byte
	O       chan<- []byte
}

func NewMqtt(params string, i <-chan []byte, o chan<- []byte) *Mqtt {
	opt, err := url.Parse(params)
	if err != nil {
		return nil
	}
	password, _ := opt.User.Password()

	conn, err := net.Dial("tcp", "localhost:1883")
	if err != nil {
		fmt.Println(err)
	}

	return &Mqtt{
		I: i,
		O: o,
		Client: paho.NewClient(paho.ClientConfig{
			ClientID: "ttttt",
			Conn:     conn,
		}),
		Connect: paho.ConnectFromPacketConnect(&packets.Connect{
			WillMessage: []byte("233"),
			Password:    []byte(password),
			Username:    opt.User.Username(),
			ClientID:    "dev",
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
	defer fmt.Println("MQTT exit")
	t.Client.Connect(ctx, t.Connect)

	for {
		select {
		case raw := <-t.I:
			fmt.Println("SSS:", string(raw))

			for topic, data := range util.DetachTran(raw) {
				res, err := t.Client.Publish(context.TODO(), &paho.Publish{
					// TODO: config
					//Payload    []byte
					Payload: data,
					//Topic      string
					Topic: "test/" + topic,
					//Properties *Properties
					//PacketID   uint16
					//QoS        byte
					//Duplicate  bool
					//Retain     bool

				})
				// TODO: error log
				if err != nil {
					//logger.Println(err)
					continue
				}
				if res != nil && res.ReasonCode != packets.PubrecSuccess {
					//logger.Println(err)
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
