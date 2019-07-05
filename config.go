package main

import (
  "io/ioutil"

  yaml "gopkg.in/yaml.v2"
)

type Server struct {
    Id int
    Secret_key string
    Mqtt string
    Tcpc string
    Tcps string
    Tran string
}

type Ncp struct {
  Download struct {
    Map string
  }
  Upload struct {
    Map string
    Air_log string
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
  Log_level string
  Ncp
}


func getConfig(str string) (Config, error) {
  config := Config{}
  configFile, err := ioutil.ReadFile(str)
  yaml.Unmarshal(configFile, &config)
  return config, err
}

