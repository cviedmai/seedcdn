package demultiplexer

import (
  "io"
  "log"
  "sync"
  "net/http"
  "seedcdn/core"
)

const (
  IDEAL_CHUNK_COUNT = 100
  CHUNKLET_SIZE = core.CHUNK_SIZE / IDEAL_CHUNK_COUNT + IDEAL_CHUNK_COUNT
)

type Payload struct {
  Header http.Header
  Status int
  Data []byte
  Finished bool
}

type Master struct {
  key string
  lock sync.Mutex
  Observers []chan *Payload
}

var (
  proxyHeaders = []string{"Content-Length", "Content-Range", "Content-Type", "Cache-Control"}
  errorPayload = &Payload{Header: make(http.Header), Status: 500, Data: []byte{}, Finished: true,}
)
func (m *Master) Observed(observer chan *Payload) {
  m.lock.Lock()
  defer m.lock.Unlock()
  m.Observers = append(m.Observers, observer)
}

func (m *Master) Run(response *http.Response, err error, masterHandler Handler) {
  if response != nil && response.Body != nil { defer response.Body.Close() }
  if err != nil {
    m.finish(errorPayload, err)
    return
  }

  status := response.StatusCode
  header := make(http.Header, len(proxyHeaders))
  for _, h := range proxyHeaders {
    value := response.Header.Get(h)
    if len(value) > 0 { header.Set(h, value) }
  }
  if response.ContentLength < 1 {
    m.finish(&Payload{header, status, nil, true,}, nil)
    return
  }

  data := make([]byte, response.ContentLength)
  read := 0
  for {
    n, err := response.Body.Read(data[read:])
    if n > 0 {
      read += n
      m.flush(&Payload{header, status, data[0:read], false,})
    }
    if err == io.EOF {
      break
    } else if err != nil {
      m.finish(errorPayload, err)
      return
    }
  }
  final := &Payload{header, status, data[0:read], true,}
  //Flush the slaves (which releases them) before we do any IO
  m.flush(final)
  masterHandler(final)
  m.finish(final, nil)
}

func (m *Master) flush(payload *Payload) {
  m.lock.Lock()
  defer m.lock.Unlock()
  for _, observer := range m.Observers {
    go func (o chan *Payload) { o <- payload }(observer)
  }
  if payload.Finished { m.Observers = make([]chan *Payload, 1) }
}

func (m *Master) finish(payload *Payload, err error) {
  if err != nil { log.Println("master: ", err) }
  Cleanup(m.key)
  //This is important to call even in a non-error case,
  //because new slaves could have joined before we cleaned up
  m.flush(payload)
}
