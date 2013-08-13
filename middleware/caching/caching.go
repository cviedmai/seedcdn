package caching

import (
  "os"
  "log"
  "path"
  "net/http"
  "encoding/gob"
  "seedcdn/core"
  "seedcdn/demultiplexer"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  //todo consistent hash around a configurable number of drives/paths
  root := "./storage"
  dir := path.Join(context.FileKey[0:2], context.FileKey[0:4], context.FileKey)
  file := context.Key
  demultiplexer.Demultiplex(context, toResponse(res), toDisk(root, dir, file))
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

func toDisk(root, dir, file string) demultiplexer.Handler {
  return func(payload *demultiplexer.Payload) {
    tmp := path.Join(root, "tmp", file)
    f, err := os.Create(tmp)
    if err != nil {
      log.Println(err)
      return
    }

    gob.NewEncoder(f).Encode(payload)
    f.Close()

    dir := path.Join(root, dir)
    if err = os.MkdirAll(dir, 0744); err != nil { log.Println(err) }
    if err = os.Rename(tmp, path.Join(dir, file)); err != nil { log.Println(err) }
  }
}
