package seedcdn

import (
  "log"
  "time"
  "runtime"
  "net/http"
  "io/ioutil"
  "seedcdn/core"
)

func Run() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  data, err := ioutil.ReadFile("config.json")
  if err != nil { log.Fatal(err) }
  core.LoadConfig(data)

  s := &http.Server {
    Addr: core.GetConfig().Listen,
    Handler: new(Handler),
    ReadTimeout: 10 * time.Second,
    MaxHeaderBytes: 8192,
  }
  log.Fatal(s.ListenAndServe())
}
