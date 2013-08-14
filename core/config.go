package core

import (
  "log"
  "io/ioutil"
  "encoding/json"
)

type Config struct{
  Listen string
  Upstream string
  Drives []string
}

var config *Config

func GetConfig() *Config {
  return config
}

func init () {
  data, err := ioutil.ReadFile("config.json")
  if err != nil { log.Fatal(err) }
  loadConfig(data)
}

func loadConfig(data []byte) {
  config = new(Config)
  if err := json.Unmarshal(data, config); err != nil {
    log.Fatal("parse config: ", err)
  }
}
