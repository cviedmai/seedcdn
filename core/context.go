package core

import (
  "fmt"
  "strconv"
  "net/http"
  "crypto/md5"
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
  bucket string

  Key string
  Dir string
  Fixed bool
  TempDir string
  Chunks []*Chunk
  HeaderFile string
  Req *http.Request
}

func NewContext(req *http.Request) *Context {
  c := &Context {Req: req,}
  c.bucket = Hash(req.URL.Path)
  drive, _ := drives.Get(c.bucket)
  c.Dir = drive + "/" + c.bucket[0:2] + "/" + c.bucket[0:4] + "/" + c.bucket + "/"
  c.TempDir = drive + "/tmp/"
  c.HeaderFile = c.Dir + c.bucket + ".hdr"
  return c
}

func (c *Context) ChunkFile(chunk *Chunk) string {
  return c.bucket + "_" + strconv.Itoa(chunk.N)
}

func Hash(value string) (string) {
  h := md5.New()
  h.Write([]byte(value))
  return fmt.Sprintf("%x", h.Sum(nil))
}
