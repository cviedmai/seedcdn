package core

import (
  "fmt"
  "math"
  "strconv"
  "net/http"
  "crypto/md5"
  "seedcdn/header"
  "github.com/stathat/consistent"
)

var drives *consistent.Consistent
func init() {
  drives = consistent.New()
  for _, drive := range GetConfig().Drives {
    drives.Add(drive)
  }
}

type Context struct {
  Key string
  Dir string
  TempDir string
  Chunks []*Chunk
  DataFile string
  HeaderFile string
  Req *http.Request
}

func NewContext(req *http.Request) *Context {
  c := &Context {Req: req,}
  bucket := Hash(req.URL.Path)
  drive, _ := drives.Get(bucket)
  c.Dir = drive + "/" + bucket[0:2] + "/" + bucket[0:4] + "/" + bucket + "/"
  c.TempDir = drive + "/tmp/"
  c.HeaderFile = c.Dir + bucket + ".hdr"

  ranges := header.ParseRange(req.Header.Get("range"))
  return c
}

func Hash(value string) (string) {
  h := md5.New()
  h.Write([]byte(value))
  return fmt.Sprintf("%x", h.Sum(nil))
}
