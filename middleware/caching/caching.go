package caching

import (
  "io"
  "os"
  "strconv"
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
  totalLength := serverChunk(context, res, first, true)
  i := (first.From / core.CHUNK_SIZE + core.CHUNK_SIZE) / core.CHUNK_SIZE
  l := totalLength / core.CHUNK_SIZE + 1
  cc := res.(http.CloseNotifier).CloseNotify()
  for ; i < l; i++ {
    select {
    case <- cc:
      return
    default:
      chunk := core.GetChunk(i)
      serverChunk(context, res, chunk, false)
    }
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
  return demultiplexer.Demultiplex(context, chunk, toResponse(context, res, chunk, first), toDisk(context))
}

func toResponse(context *core.Context, res http.ResponseWriter, chunk *core.Chunk, first bool) demultiplexer.Handler {
  sentHeaders := !first
  to := context.Range.To
  from := context.Range.From
  offset := chunk.N * core.CHUNK_SIZE
  var start int
  if first { start = from }
  return func(payload *demultiplexer.Payload) {
    if payload.Finished { return }
    if sentHeaders == false {
      rh := res.Header()
      for k, v := range payload.Header { rh[k] = v }

      if to == 0 {
        to = payload.TotalLength
        context.Range.To = payload.TotalLength
      }
      rh.Set("Content-Length", strconv.Itoa(to - from))
      rh.Set("Accept-Ranges", "bytes")
      if context.Range.RangeRequest {
        rh.Set("Content-Range", "bytes " + strconv.Itoa(from) + "-" + strconv.Itoa(to-1) + "/" + strconv.Itoa(payload.TotalLength))
        res.WriteHeader(206)
      } else {
        res.WriteHeader(200)
      }
      sentHeaders = true
    }
    l := len(payload.Data)
    end := l
    if (offset + end) > to { end = to - offset }
    if start >= end { return }
    println(chunk.N, start, end)
    res.Write(payload.Data[start:end])
    start = l + 1
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
