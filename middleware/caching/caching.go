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
  "github.com/stathat/consistent"
)

type CacheHeader struct {
  Header http.Header
  Status int
}

var pool = bytepool.New(1024, 2048)
var drives *consistent.Consistent

func init() {
  drives = consistent.New()
  for _, drive := range core.GetConfig().Drives {
    drives.Add(drive)
  }
}

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  //todo consistent hash around a configurable number of drives/paths
  root, _ := drives.Get(context.Bucket)
  if fromDisk(root, res, context) {
    core.Stats.CacheHit()
    return
  }
  core.Stats.CacheMiss()
  demultiplexer.Demultiplex(context, toResponse(res), toDisk(root, context))
}


func fromDisk(root string, res http.ResponseWriter, context *core.Context) bool {
  headerFile, err := os.Open(root + context.HeaderFile)
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
  res.Header().Set("X-Accel-Redirect", root + context.DataFile)
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

func toDisk(root string, context *core.Context) demultiplexer.Handler {
  return func(payload *demultiplexer.Payload) {
    if err := os.MkdirAll(root + context.Dir, 0744); err != nil {
      log.Println("mkdir: ", err)
      return
    }
    if write(root, context.Key, context.DataFile, payload.Data) {
      bytes := pool.Checkout()
      err := gob.NewEncoder(bytes).Encode(&CacheHeader{payload.Header, payload.Status})
      if err != nil { println(err.Error()) }
      write(root, context.Key, context.HeaderFile, bytes.Bytes())
      bytes.Close()
    }
  }
}

func write(root, key, file string, data []byte) bool {
  tmp := path.Join(root, "tmp", key)
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
  if err = os.Rename(tmp, root + file); err != nil {
    log.Println("rename: ", err)
    return false
  }
  return true
}
