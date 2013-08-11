package core

import (
  "net/http"
)

type Middleware func (context *Context, res http.ResponseWriter)

type MiddlewareLink func(context *Context, res http.ResponseWriter, next Middleware)

type MiddlewareWrapper struct {
  Middleware MiddlewareLink
  Next *MiddlewareWrapper
}

func (wrapper *MiddlewareWrapper) Yield(context *Context, res http.ResponseWriter) {
  var next  Middleware
  if wrapper.Next != nil {
    next = wrapper.Next.Yield
  }
  wrapper.Middleware(context, res, next)
}
