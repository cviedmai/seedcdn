package core

import (
  "log"
  "io/ioutil"
  "encoding/json"
)

type Config struct{
  Listen string
}

var config = loadConfig()

func GetConfig() *Config {
  return config
}

func loadConfig() *Config{
  c := new(Config)
  data, err := ioutil.ReadFile("config.json")
  if err != nil { log.Fatal(err) }
  if json.Unmarshal(data, c); err != nil {
    log.Fatal("parse config: ", err)
  }
  return c
}
