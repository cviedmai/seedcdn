package demultiplexer

import (
  "net/http"
)

var (
  proxyHeaders = []string{"Content-Length", "Content-Range", "Content-Type", "Cache-Control"}
)

type Master struct {
  key string
  status int
  header http.Header
}

func (m *Master) Run(response http.Response) {
  m.status = response.StatusCode
  for _, h := range proxyHeaders {
    m.header[h] = response.Header[h]
  }
}
