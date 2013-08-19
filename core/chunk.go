package core

import (
  "strconv"
  "net/http"
)

var chunks map[int]*Chunk

type Chunk struct {
  From int
  To int
  From64 int64
  To64 int64
  N int
  Header http.Header
}

func init() {
  chunks = make(map[int]*Chunk, 10000)
  for i := 0; i < 10000; i++ {
    chunks[i] = makeChunk(i)
  }
}

func GetChunk(n int) *Chunk {
  return chunks[n]
}

func makeChunk(i int) *Chunk {
  from := i * CHUNK_SIZE
  to := (i+1) * CHUNK_SIZE - 1

  header := make(http.Header, 1)
  header.Set("Range", "bytes=" + strconv.Itoa(from) + "-" + strconv.Itoa(to))

  return &Chunk{
    From: from,
    To: to,
    From64: int64(from),
    To64: int64(to),
    N: i,
    Header: header,
  }
}
