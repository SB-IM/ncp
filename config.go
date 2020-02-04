package main

import (
  "io/ioutil"

	//"github.com/pion/webrtc/v2"
  yaml "gopkg.in/yaml.v2"
)

type Server struct {
    Id int
    Link_id int
    Secret_key string
    Mqtt string
    Tcpc string
    Tcps string
    Tran string
}

type Ncp struct {
  Common struct {
    Id string
    SecretKey string
  }
  Download map[string]string
  Upload map[string]string
  Webrtc struct {
		//Iceserver []webrtc.ICEServer
		Driver string
    Args string
  }
  Status
  Shell struct {
    Path string
    Prefix string
    Suffix string
  }
}

type Status struct {
  Link_id int `json:"link_id"`
  Position_ok bool `json:"position_ok"`
  Lat string `json:"lat"`
  Lng string `json:"lng"`
  Alt string `json:"alt"`
}

type Config struct {
	Server
	Env string
	Log ConfigLog
	Ncp
}

func getConfig(str string) (Config, error) {
  config := Config{}
  configFile, err := ioutil.ReadFile(str)
  yaml.Unmarshal(configFile, &config)
  return config, err
}

