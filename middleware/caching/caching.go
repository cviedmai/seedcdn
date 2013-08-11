package caching

import (
  "net/http"
  "seedcdn/core"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  next(context, res)
}
