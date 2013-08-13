package core

import (
  "fmt"
  "math"
  "strconv"
  "net/http"
  "crypto/md5"
  "seedcdn/header"
)

const CHUNK_SIZE = 2*1024*1024

type Context struct {
  Chunk int
  Key string
  FileKey string
  Req *http.Request
}

func NewContext(req *http.Request) *Context {
  c := &Context {
    Req: req,
  }
  r := header.ParseRange(req.Header.Get("range"))
  if len(r) == 0 {
    c.Chunk = 0
  } else {
    c.Chunk = int(math.Floor(float64(r[0].From) / float64(CHUNK_SIZE)))
  }
  c.FileKey = Hash(req.URL.Path)
  c.Key = c.FileKey + "_" + strconv.Itoa(c.Chunk)
  return c
}

func Hash(value string) (string) {
  h := md5.New()
  h.Write([]byte(value))
  return fmt.Sprintf("%x", h.Sum(nil))
}
