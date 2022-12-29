package help

import (
	"os"
	"os/exec"
	"strings"

	"sb.im/ncp/ncpio"
	"sb.im/ncp/util"
)

type Config struct {
	NcpIO []ncpio.Config `json:"ncpio"`
}

func GetConfig(str string) (Config, error) {
	config := Config{}
	return config, util.LoadConfig(str, &config)
}

func GenerateMqttConfig(name string, config []byte) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	if _, err := file.Write(config); err != nil {
		panic(err)
	}
}

func CmdRun(str string) ([]byte, error) {
	cmdArr := strings.Split(str, " ")
	return exec.Command(cmdArr[0], cmdArr[1:]...).CombinedOutput()
}
