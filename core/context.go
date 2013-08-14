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
  Chunk int
  Key string
  Dir string
  TempDir string
  DataFile string
  HeaderFile string
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
  bucket := Hash(req.URL.Path)
  drive, _ := drives.Get(bucket)
  c.Key = bucket + "_" + strconv.Itoa(c.Chunk)
  c.Dir = drive + "/" + bucket[0:2] + "/" + bucket[0:4] + "/" + bucket + "/"
  c.TempDir = drive + "/tmp/"
  c.DataFile = c.Dir + c.Key + ".dat"
  c.HeaderFile = c.Dir + c.Key + ".hdr"
  return c
}

func Hash(value string) (string) {
  h := md5.New()
  h.Write([]byte(value))
  return fmt.Sprintf("%x", h.Sum(nil))
}
