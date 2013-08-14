package caching

import (
  "os"
  "log"
  "path"
  "net/http"
  "encoding/gob"
  "seedcdn/core"
  "seedcdn/demultiplexer"
  "github.com/viki-org/bytepool"
)

type CacheHeader struct {
  Header http.Header
  Status int
}

var pool = bytepool.New(1024, 2048)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  if fromDisk(res, context) {
    core.Stats.CacheHit()
    return
  }
  core.Stats.CacheMiss()
  demultiplexer.Demultiplex(context, toResponse(res), toDisk(context))
}


func fromDisk(res http.ResponseWriter, context *core.Context) bool {
  headerFile, err := os.Open(context.HeaderFile)
  if err != nil { return false }
  defer headerFile.Close()

  ch := new(CacheHeader)
  if err := gob.NewDecoder(headerFile).Decode(ch); err != nil {
    log.Println("header decode: ", err)
    return false
  }
  for key, value := range ch.Header {
    res.Header()[key] = value
  }
  res.Header().Set("X-Accel-Redirect", context.DataFile)
  res.WriteHeader(ch.Status)
  return true
}

func toResponse(res http.ResponseWriter) demultiplexer.Handler {
  var read int
  var sentHeaders bool
  return func(payload *demultiplexer.Payload) {
    if sentHeaders == false {
      for key, value := range payload.Header {
        res.Header()[key] = value
      }
      res.Header().Set("Accept-Ranges", "bytes")
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

func toDisk(context *core.Context) demultiplexer.Handler {
  return func(payload *demultiplexer.Payload) {
    if err := os.MkdirAll(context.Dir, 0755); err != nil {
      log.Println("mkdir: ", err)
      return
    }
    if write(context.TempDir, context.Key, context.DataFile, payload.Data) {
      bytes := pool.Checkout()
      err := gob.NewEncoder(bytes).Encode(&CacheHeader{payload.Header, payload.Status})
      if err != nil { println(err.Error()) }
      write(context.TempDir, context.Key, context.HeaderFile, bytes.Bytes())
      bytes.Close()
    }
  }
}

func write(tempDir, key, file string, data []byte) bool {
  tmp := path.Join(tempDir + key)
  f, err := os.Create(tmp)
  if err != nil {
    log.Println("create tmp: ", err)
    return false
  }
  defer f.Close()
  _, err = f.Write(data)
  if err != nil {
    log.Println("write: ", err)
    return false
  }
  if err = os.Rename(tmp, file); err != nil {
    log.Println("rename: ", err)
    return false
  }
  return true
}
