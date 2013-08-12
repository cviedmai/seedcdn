package demultiplexer

import (
  "io"
  "time"
  "net/http"
  "seedcdn/core"
)

const (
  IDEAL_CHUNK_COUNT = 10
  CHUNKLET_SIZE = core.CHUNK_SIZE / IDEAL_CHUNK_COUNT + IDEAL_CHUNK_COUNT
)

var proxyHeaders = []string{"Content-Length", "Content-Range", "Content-Type", "Cache-Control"}

type Sync interface {
  Header() http.Header
  Data() []byte
  Status() int
}

type Master struct {
  read int
  key string
  status int
  data []byte
  sync chan Sync
  header http.Header
  observers chan chan []byte
}

func New(key string) *Master{
  return &Master{
    key: key,
    header: make(http.Header, len(proxyHeaders)),
    observers: make(chan chan []byte, 16),
  }
}

func (m *Master) Run(response http.Response) {
  //handle body close, and errors

  m.status = response.StatusCode
  for _, h := range proxyHeaders {
    m.header[h] = response.Header[h]
  }
  //todo pull this from a pool
  m.data = make([]byte, core.CHUNK_SIZE)
  for {
    m.Sync()
    n, err := response.Body.Read(m.data[m.read:m.read+CHUNKLET_SIZE])
    if n > 0 { m.Notify(m.data[m.read:m.read+n]) }
    if err == io.EOF { break }
    if err != nil {
      //todo
    }
    m.read += n
  }
}

func (m *Master) Header() http.Header {
  return m.header
}

func (m *Master) Status() int {
  return m.status
}

func (m *Master) Data() []byte {
  return m.data[0:m.read]
}

func (m *Master) Sync() {
  for {
    select {
    case m.sync <- m:
    default:
      break
    }
  }
}

func (m *Master) Notify(data []byte) {
  for observer := range m.observers {
    select {
    case observer <- data:
    case <- time.After(time.Second * 5): // todo configure this, but what's a good value?!
      continue
    }
  }
}
