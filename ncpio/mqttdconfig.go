package ncpio

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
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

	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return config.Mqttd, err
	} else {
		err = yaml.Unmarshal(configFile, &config)
		return config.Mqttd, err
	}
}
