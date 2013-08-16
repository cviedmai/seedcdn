package seedcdn

import (
  "net/http"
  "seedcdn/core"
  "seedcdn/middleware/logs"
  "seedcdn/middleware/purge"
  "seedcdn/middleware/chunker" 
  "seedcdn/middleware/caching"
)

type Handler struct{}
var head *core.MiddlewareWrapper

func init() {
  head = &core.MiddlewareWrapper {Middleware: purge.Run}
  prev := head
  for _, middleware := range []core.MiddlewareLink{logs.Run, chunker.Run, caching.Run} {
    wrapper := &core.MiddlewareWrapper {Middleware: middleware,}
    prev.Next = wrapper
    prev = wrapper
  }
}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  core.Stats.Request()
  context := core.NewContext(req)
  head.Yield(context, res)
}