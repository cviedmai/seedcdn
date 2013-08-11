package core

import (
  "math"
  "net/http"
  "seedcdn/header"
)

const CHUNKSIZE = float64(2*1024*1024)

type Context struct {
  Req *http.Request
  Chunk int
}

func NewContext(req *http.Request) *Context {
  c := &Context {
    Req: req,
  }
  r := header.ParseRange(req.Header.Get("range"))
  if len(r) == 0 {
    c.Chunk = 0
  } else {
    c.Chunk = int(math.Floor(float64(r[0].From) / CHUNKSIZE))
  }
  return c
}
