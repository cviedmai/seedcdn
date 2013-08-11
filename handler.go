package seedcdn

import (
  "net/http"
  "seedcdn/core"
  "seedcdn/middleware/proxy"
  "seedcdn/middleware/caching"
)

type Handler struct{}
var head *core.MiddlewareWrapper

func init() {
  head = &core.MiddlewareWrapper {Middleware: caching.Run}
  prev := head
  for _, middleware := range []core.MiddlewareLink{proxy.Run} {
    wrapper := &core.MiddlewareWrapper {Middleware: middleware,}
    prev.Next = wrapper
    prev = wrapper
  }
}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  context := core.NewContext(req)
  head.Yield(context)
}
