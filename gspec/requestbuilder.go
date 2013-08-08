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

func (rb *RequestBuilder) Request() *http.Request {
  return rb.req
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {

}
