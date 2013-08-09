package gspec

import (
  "net/http"
)

type RequestBuilder struct {
  req *http.Request
}

func Request() *RequestBuilder {
  return &RequestBuilder{
    req: new(http.Request),
  }
}

func (rb *RequestBuilder) Req() *http.Request {
  return rb.req
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
  rb.req.Header.Set(key, value)
  return rb
}
