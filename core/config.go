package core

import (
  "log"
  "encoding/json"
)

type Config struct{
  Listen string
  Upstream string
}

var config *Config

func GetConfig() *Config {
  return config
}

func LoadConfig(data []byte) {
  config = new(Config)
  if err := json.Unmarshal(data, config); err != nil {
    log.Fatal("parse config: ", err)
  }
}
