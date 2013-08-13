package seedcdn

import (
  "net/http"
  "seedcdn/core"
  "seedcdn/middleware/caching"
)

type Handler struct{}
var head *core.MiddlewareWrapper

func init() {
  head = &core.MiddlewareWrapper {Middleware: caching.Run}
}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  context := core.NewContext(req)
  head.Yield(context, res)
}
