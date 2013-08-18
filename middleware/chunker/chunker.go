package chunker

import (
  "net/http"
  "seedcdn/core"
  "seedcdn/header"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  ranges := header.ParseRange(context.Req.Header.Get("range"))
  context.Chunks, context.Fixed = calculateChunks(context, ranges[0])
  next(context, res)
}

func calculateChunks(context *core.Context, r header.Range) ([]*core.Chunk, bool) {
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
