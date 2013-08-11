package core

type Middleware func (context *Context)

type MiddlewareLink func(context *Context, next Middleware)

type MiddlewareWrapper struct {
  Middleware MiddlewareLink
  Next *MiddlewareWrapper
}

func (wrapper *MiddlewareWrapper) Yield(context *Context) {
  var next  Middleware
  if wrapper.Next != nil {
    next = wrapper.Next.Yield
  }
  wrapper.Middleware(context, next)
}
