package chunker

import (
  "math"
  "strconv"
  "net/http"
  "seedcdn/core"
  "seedcdn/header"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  ranges := header.ParseRange(context.Req.Header.Get("range"))
  context.Chunks = calculateChunks(context, ranges[0])
  next(context, res)
}


func calculateChunks(context *core.Context, r header.Range) []core.Chunk {
  if r.To == 0 {
    return []core.Chunk{*&core.Chunk{From: r.From,}}
  }

  chunks := make([]core.Chunk, 0, 2)
  for i := r.From; i <= r.To; i += core.CHUNK_SIZE {
    n := int(math.Floor(float64(i) / float64(core.CHUNK_SIZE)))
    from := i
    if i != r.From { from = core.CHUNK_SIZE * n }

    to := (n + 1) * core.CHUNK_SIZE - 1
    if to > r.To { to = r.To }
    key := context.File("_" + strconv.Itoa(n))

    chunk := &core.Chunk{
      From: from,
      To: to,
      Key: key,
      DataFile: context.Dir + key,
    }
    chunks = append(chunks, *chunk)
  }
  return chunks
}
