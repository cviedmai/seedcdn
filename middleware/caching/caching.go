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
  root := "/"
  dir := path.Join(root, context.FileKey[0:2], context.FileKey[0:4], context.FileKey)
  fullPath := path.Join(dir, context.FileKey)
  if file, err := os.Open(fullPath); err == nil {
    fromFile(res, file)
    file.Close()
    return
  }
  demultiplexer.Demultiplex(context, toResponse(res), toDisk(root, dir, fullPath, context.FileKey))
}

func fromFile(res http.ResponseWriter, file *os.File) {
  payload := new(demultiplexer.Payload)
  if err := gob.NewDecoder(file).Decode(payload); err != nil {
    log.Println("gob decode: ", err)
  }
  for key, value := range payload.Header {
    res.Header()[key] = value
  }
  res.WriteHeader(payload.Status)
  res.Write(payload.Data)
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

func toDisk(root, dir, fullPath, file string) demultiplexer.Handler {
  return func(payload *demultiplexer.Payload) {
    tmp := path.Join(root, "tmp", file)
    f, err := os.Create(tmp)
    if err != nil {
      log.Println("create tmp:", err)
      return
    }

    gob.NewEncoder(f).Encode(payload)
    f.Close()

    if err = os.MkdirAll(dir, 0744); err != nil { log.Println("mkdir: ", err) }
    if err = os.Rename(tmp, fullPath); err != nil { log.Println("rename: ", err) }
  }
}
