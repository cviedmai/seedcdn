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
  for _, chunk := range context.Chunks {
    if chunk.To == 0 {
      //todo get entire file
    } else {
      serverChunk(context, res, &chunk)
    }
  }
}

func serverChunk(context *core.Context, res http.ResponseWriter, chunk *core.Chunk) {
  if file, err := os.Open(chunk.DataFile); err == nil {
    core.Stats.CacheHit()
    if chunk.From > 0 { file.Seek(chunk.From) }
    io.CopyN(res, file, chunk.To - chunk.From)
    file.Close()
  } else {
    core.Stats.CacheMiss()
    demultiplexer.Demultiplex(context, toResponse(res, chunk), toDisk(context))
}


func toResponse(res http.ResponseWriter, chunk *core.Chunk) demultiplexer.Handler {
  from := chunk.From
  return func(payload *demultiplexer.Payload) {
    to := len(payload.Data)
    if to > from
      if to > chunk.To { to = chunk.To }
      res.Write(payload.Data[read:to])
      from = to
    }
  }
}

func toDisk(context *core.Context) demultiplexer.Handler {
  return func(payload *demultiplexer.Payload) {
    // if err := os.MkdirAll(context.Dir, 0755); err != nil {
    //   log.Println("mkdir: ", err)
    //   return
    // }
    // if write(context.TempDir, context.Key, context.DataFile, payload.Data) {
    //   bytes := pool.Checkout()
    //   err := gob.NewEncoder(bytes).Encode(&CacheHeader{payload.Header, payload.Status})
    //   if err != nil { println(err.Error()) }
    //   write(context.TempDir, context.Key, context.HeaderFile, bytes.Bytes())
    //   bytes.Close()
    // }
  }
}

func write(tempDir, key, file string, data []byte) bool {
  // tmp := path.Join(tempDir + key)
  // f, err := os.Create(tmp)
  // if err != nil {
  //   log.Println("create tmp: ", err)
  //   return false
  // }
  // defer f.Close()
  // _, err = f.Write(data)
  // if err != nil {
  //   log.Println("write: ", err)
  //   return false
  // }
  // if err = os.Rename(tmp, file); err != nil {
  //   log.Println("rename: ", err)
  //   return false
  // }
  // return true
}



// func fromDisk(res http.ResponseWriter, context *core.Context) bool {
//   headerFile, err := os.Open(context.HeaderFile)
//   if err != nil { return false }
//   defer headerFile.Close()

//   ch := new(CacheHeader)
//   if err := gob.NewDecoder(headerFile).Decode(ch); err != nil {
//     log.Println("header decode: ", err)
//     return false
//   }
//   for key, value := range ch.Header {
//     res.Header()[key] = value
//   }
//   res.Header().Set("X-Accel-Redirect", context.DataFile)
//   res.WriteHeader(ch.Status)
//   return true
// }
