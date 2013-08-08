package seedcdn

import (
  "math"
  "net/http"
)

const CHUNKSIZE = 2*1024*1024

type Context struct {
  req *http.Request
  chunk int
}

func NewContext(req *http.Request) *Context {
  c := &Context {
    req: req,
  }
  r := header.ParseRange(req.Header.Get("range"))
  if len(r) == 0 {
    c.chunk = 0
  } else {
    c.chunk = math.Floor(r[0].from / CHUNKSIZE)
  }
  return c
}
