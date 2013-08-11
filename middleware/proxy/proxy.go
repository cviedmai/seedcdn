package proxy

import (
  "seedcdn/core"
)

func Run (context *core.Context, next core.Middleware) {
  println("proxy")
  return
}
