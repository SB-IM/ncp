package integration

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
	"sb.im/ncp/ncpio"
)

type Config struct {
	NcpIO []ncpio.Config `json:"ncpio"`
}

func getConfig(str string) (Config, error) {
	config := Config{}
	configFile, err := ioutil.ReadFile(str)
	if err != nil {
		return config, err
	} else {
		err = yaml.Unmarshal(configFile, &config)
		return config, err
	}
}

func generateMqttConfig(name string, config []byte) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	if _, err := file.Write(config); err != nil {
		panic(err)
	}
}
