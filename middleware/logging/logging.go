package logging

import (
  "seedcdn/core"
)

func Run (context *core.Context, next core.Middleware) {
  println("logging")
  next(context)
}
