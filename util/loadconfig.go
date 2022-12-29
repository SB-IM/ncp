package util

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func LoadConfig(path string, config interface{}) error {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(os.ExpandEnv(string(configFile))), config)
}
