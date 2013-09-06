package chunker

import (
  "net/http"
  "seedcdn/core"
  "seedcdn/header"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  r := context.Req.Header.Get("range")
  if len(r) == 0 {
    context.Range = &header.Range{0, 0, false}
  } else {
    context.Range = header.ParseRange(r)[0]
  }
  context.Chunks, context.Fixed = calculateChunks(context)
  next(context, res)
}

func calculateChunks(context *core.Context) ([]*core.Chunk, bool) {
  r := context.Range
  if r.To == 0 {
    n := r.From / core.CHUNK_SIZE
    return []*core.Chunk{core.GetChunk(n)}, false
  }

  chunks := make([]*core.Chunk, 0, 2)
  for i := r.From; i <= r.To; i += core.CHUNK_SIZE {
    n := i / core.CHUNK_SIZE
    chunks = append(chunks, core.GetChunk(n))
  }
  return chunks, true
}
