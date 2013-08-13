package caching

import (
  "net/http"
  "seedcdn/core"
  "seedcdn/demultiplexer"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  demultiplexer.Demultiplex(context, toResponse(res), toDisk())
}

func toResponse(res http.ResponseWriter) demultiplexer.Handler {
  var read int
  var sentHeaders bool
  return func(payload *demultiplexer.Payload) {
    if sentHeaders == false {
      for key, value := range payload.Header {
        res.Header()[key] = value
      }
      res.WriteHeader(payload.Status)
      sentHeaders = true
    }
    length := len(payload.Data)
    if length > read {
      res.Write(payload.Data[read:])
      read = length
    }
  }
}

func toDisk() demultiplexer.Handler {
  return func(payload *demultiplexer.Payload) {
    if payload.Finished == false { return }
  }
}
