package caching

import (
  "seedcdn/core"
)

func Run (context *core.Context, next core.Middleware) {
  next(context)
}
