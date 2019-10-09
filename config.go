package main

import (
  "io/ioutil"
  "fmt"

  yaml "gopkg.in/yaml.v2"
)

type Server struct {
    Id int
    Secret_key string
    Mqtt string
    Tcpc string
    Tcps string
}

type Ncp struct {
  Download map[string]string
  Upload map[string]string
  Live struct {
    Args string
  }
  Status struct {
    Link_id int
    Position_ok bool
    Lat string
    Lng string
    Alt string
  }
  Shell struct {
    Path string
    Prefix string
    Suffix string
  }
}

func (this *Ncp) Upload2 () {
  fmt.Println((*this).Upload["map"])
  fmt.Println((*this).Upload["map2"])
  if (*this).Upload["map2"] == "" {
    fmt.Println("EEEEEEEEEEEEE")
  }
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

