package demultiplexer

import (
  "io"
  "sync"
  "net/http"
  "seedcdn/core"
)

const (
  IDEAL_CHUNK_COUNT = 10
  CHUNKLET_SIZE = core.CHUNK_SIZE / IDEAL_CHUNK_COUNT + IDEAL_CHUNK_COUNT
)

var proxyHeaders = []string{"Content-Length", "Content-Range", "Content-Type", "Cache-Control"}

type Payload struct {
  Header http.Header
  Data []byte
  Status int
  Finished bool
}

type Master struct {
  key string
  lock sync.Mutex
  Observers []chan *Payload
}

func New(key string) *Master{
  return &Master{key: key,}
}

func (m *Master) Observed(observer chan *Payload) {
  m.lock.Lock()
  defer m.lock.Unlock()
  m.Observers = append(m.Observers, observer)
}

func (m *Master) Run(response *http.Response, err error, masterHandler Handler) {
  //todo handle errors
  defer response.Body.Close()

  status := response.StatusCode
  header := make(http.Header, len(proxyHeaders))
  for _, h := range proxyHeaders {
    value := response.Header.Get(h)
    if len(value) > 0 { header.Set(h, value) }
  }

  data := make([]byte, core.CHUNK_SIZE)
  read := 0
  for {
    n, err := response.Body.Read(data[read:read+CHUNKLET_SIZE])
    if n > 0 {
      read += n
      m.flush(&Payload{Header: header, Status: status, Data: data[0:read], Finished: false,})
    }
    if err == io.EOF {
      break
    } else if err != nil {
      //todo
    }
  }
  final := &Payload{Header: header, Status: status, Data: data[0:read], Finished: true,}
  //Flush the slaves (which releases them) before we do any IO
  m.flush(final)
  masterHandler(final)

  //todo, linger around a bit to fufill any late joiners?
  Cleanup(m.key)
  //Maybe some new slaves joined before we cleaned up
  m.flush(final)
}

func (m *Master) flush(payload *Payload) {
  m.lock.Lock()
  defer m.lock.Unlock()
  for _, observer := range m.Observers {
    go func (o chan *Payload) { o <- payload }(observer)
  }
  if payload.Finished { m.Observers = make([]chan *Payload, 1) }
}
