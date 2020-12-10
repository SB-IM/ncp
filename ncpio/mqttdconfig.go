package ncpio

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type MqttdConfig struct {
	ID      int
	Client  string
	Status  string
	Network string
	Broker  string
	Static  *StaticStatus
	Rpc     struct {
		I string
		O string
	}
	Gtran struct {
		Prefix string
	}
	Trans map[string]struct {
		Retain bool
		QoS    byte
	}
}

func loadMqttConfigFromFile(file string) (*MqttdConfig, error) {
	config := struct {
		Mqttd *MqttdConfig
	}{}

	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return config.Mqttd, err
	} else {
		err = yaml.Unmarshal(configFile, &config)
		return config.Mqttd, err
	}
}
