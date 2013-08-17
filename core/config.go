package core

import (
  "os"
  "log"
  "strings"
  "io/ioutil"
  "encoding/json"
)

const CHUNK_SIZE = 5*1024*1024

type Config struct{
  RedisDB int
  Redis string
  Listen string
  Prefix string
  Upstream string
  Drives []string
  RedisProtocol string
  PurgeWhiteList map[string]bool
  RangedExtensions map[string]bool
}

var config *Config

func GetConfig() *Config {
  return config
}

func init () {
  //this is sooo horrible
  if strings.Contains(os.Args[0], ".test")  {
    config = &Config{
      Upstream: "test.viki.io",
      Drives: []string{"/tmp"},
      PurgeWhiteList: map[string]bool{"127.0.0.1": true},
      RangedExtensions: map[string]bool{".mp4": true},
      Redis: "localhost:6379",
      RedisDB: 5,
      RedisProtocol: "tcp",
    }
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
  if len(config.RedisProtocol) == 0 {
    config.RedisProtocol = "tcp"
    if strings.Contains(config.Redis, "sock") {
      config.RedisProtocol = "unix"
    }
  }
}
