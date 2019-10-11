package main

import (
  "io/ioutil"
  "fmt"
  "errors"

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
  Common struct {
    Id string
    SecretKey string
  }
  Download map[string]string
  Upload map[string]string
  Live struct {
    Args string
  }
  Status
  Shell struct {
    Path string
    Prefix string
    Suffix string
  }
}

type NcpCmd struct {
  config Ncp
}

func (this *NcpCmd) Download (filename, source string) error {
  if (*this).config.Download[filename] == "" {
    fmt.Println("EEEEEEEEEEEEE")
    return errors.New("No " + filename + " config found")
  } else {
    return httpDownload((*this).config.Common.Id, (*this).config.Common.SecretKey, (*this).config.Download[filename], source)
  }
}

func (this *NcpCmd) Upload (filename, target string) error {
  if (*this).config.Upload[filename] == "" {
    fmt.Println("EEEEEEEEEEEEE")
    return errors.New("No " + filename + " config found")
  } else {
    return httpUpload((*this).config.Common.Id, (*this).config.Common.SecretKey, filename, (*this).config.Upload[filename], target)
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

