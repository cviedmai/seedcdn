package core

import (
  "os"
  "log"
  "strings"
  "io/ioutil"
  "encoding/json"
)

const CHUNK_SIZE = 2*1024*1024

type Config struct{
  Listen string
  Upstream string
  Drives []string
  PurgeWhiteList map[string]bool
}

var config *Config

func GetConfig() *Config {
  return config
}

func init () {
  //this is sooo horrible
  if strings.Contains(os.Args[0], ".test")  {
    config = &Config{Upstream: "test.viki.io", Drives: []string{"/tmp"}, PurgeWhiteList: map[string]bool{"127.0.0.1": true},}
  } else {
    data, err := ioutil.ReadFile("config.json")
    if err != nil { log.Fatal(err) }
    loadConfig(data)
  }
}

func loadConfig(data []byte) {
  config = new(Config)
  if err := json.Unmarshal(data, config); err != nil {
    log.Fatal("parse config: ", err)
  }
}
