package caching

import (
  "io"
  "os"
  "net/http"
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
  if context.Fixed {
    for i, chunk := range context.Chunks {
      serverChunk(context, res, chunk, i == 0)
    }
  } else {
    drain(context, res)
  }
}

func drain(context *core.Context, res http.ResponseWriter) {
  first := context.Chunks[0]
  contentLength := serverChunk(context, res, first, true)

  i := (first.From / core.CHUNK_SIZE + core.CHUNK_SIZE) / core.CHUNK_SIZE
  l := contentLength / core.CHUNK_SIZE + 1
  for ; i < l; i++ {
    chunk := core.GetChunk(i)
    serverChunk(context, res, chunk, false)
  }
}

func serverChunk(context *core.Context, res http.ResponseWriter, chunk *core.Chunk, first bool) int {
  dataFile := context.Dir + context.ChunkFile(chunk)
  if file, err := os.Open(dataFile); err == nil {
    core.Stats.CacheHit()
    if chunk.From > 0 { file.Seek(chunk.From64, 0) }
    io.CopyN(res, file, chunk.To64 - chunk.From64)
    file.Close()
    //when called from dain, we need to get the content length!
    return 0
  }
  core.Stats.CacheMiss()
  return demultiplexer.Demultiplex(context, chunk, toResponse(res, chunk, first), toDisk(context))
}

func toResponse(res http.ResponseWriter, chunk *core.Chunk, first bool) demultiplexer.Handler {
  from := chunk.From
  sentHeaders := !first
  return func(payload *demultiplexer.Payload) {
    if sentHeaders == false {
      for k, v := range payload.Header {
        res.Header()[k] = v
      }
      sentHeaders = true
    }
    to := len(payload.Data)
    if to > from {
      if to > chunk.To { to = chunk.To }
      res.Write(payload.Data[from:to])
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
  return true
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
