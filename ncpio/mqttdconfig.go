package ncpio

import (
	"sb.im/ncp/util"
)

type MqttdConfig struct {
	ID      string
	Client  string
	Status  string
	Network string
	Broker  string
	Static  *StaticStatus
	Rpc     RpcIO
	Gtran   struct {
		Prefix string
	}
	Trans map[string]struct {
		Retain bool
		QoS    byte
	}
}

type RpcIO struct {
	QoS byte
	LRU int
	I   string
	O   string
}

func loadMqttConfigFromFile(file string) (*MqttdConfig, error) {
	config := struct {
		Mqttd *MqttdConfig
	}{
		Mqttd: &MqttdConfig{
			Rpc: RpcIO{
				QoS: 0,
				LRU: 128,
			},
		},
	}

	return config.Mqttd, util.LoadConfig(file, &config)
}
