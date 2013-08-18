package core

import (
  "strconv"
  "net/http"
)

var fixdChunks map[int]*Chunk
var openChunks map[int]*Chunk

type Chunk struct {
  From int
  To int
  From64 int64
  To64 int64
  N int
  Header http.Header
}

func init() {
  fixdChunks = make(map[int]*Chunk, 10000)
  openChunks = make(map[int]*Chunk, 10000)
  for i := 0; i < 10000; i++ {
    fixdChunks[i] = makeChunk(i, true)
    openChunks[i] = makeChunk(i, false)
  }
}

func GetChunk(n int, fixed bool) *Chunk {
  if fixed {
    return fixdChunks[n]
  }
  return openChunks[n]
}

func makeChunk(i int, fixed bool) *Chunk {
  from := i * CHUNK_SIZE
  to := 0

  var header http.Header
  if fixed {
    to = (i+1) * CHUNK_SIZE - 1
    header.Set("Range", "bytes=" + strconv.Itoa(from) + "-" + strconv.Itoa(to))
  }
  return &Chunk{
    From: from,
    To: to,
    From64: int64(from),
    To64: int64(to),
    N: i,
    Header: header,
  }
}
