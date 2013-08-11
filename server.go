package seedcdn

import (
  "time"
  "runtime"
  "net/http"
  "seedcdn/core"
)

func Run() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  s := &http.Server {
    Addr: core.GetConfig().Listen,
    Handler: new(Handler),
    ReadTimeout: 10 * time.Second,
    MaxHeaderBytes: 8192,
  }
  s.ListenAndServe()
}
