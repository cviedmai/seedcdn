package chunker

import (
  "testing"
  "net/http"
  "seedcdn/core"
  "net/http/httptest"
  "github.com/viki-org/gspec"
)

func TestHandlesAFullRange(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(context.Fixed).ToEqual(false)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  assertChunk(spec, context.Chunks[0], 0, 5242879)
}

func TestHandlesARangeWithNoLimit(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(context.Fixed).ToEqual(false)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  assertChunk(spec, context.Chunks[0], 0, 5242879)
}

func TestHandlesARangeWithNoLimitAfterTheFirstChunk(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 7000000-").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(context.Fixed).ToEqual(false)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  assertChunk(spec, context.Chunks[0], 5242880, 10485759)
}

func TestHandlesARangeWithinASingleChunk(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-2000000").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(context.Fixed).ToEqual(true)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  assertChunk(spec, context.Chunks[0], 0, 5242879)
}

func TestHandlesARangeAcrossMultipleChunks(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-12000000").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(context.Fixed).ToEqual(true)
  spec.Expect(len(context.Chunks)).ToEqual(3)
  assertChunk(spec, context.Chunks[0], 0, 5242879)
  assertChunk(spec, context.Chunks[1], 5242880, 10485759)
  assertChunk(spec, context.Chunks[2], 10485760, 15728639)
}

func TestCallsTheNextMiddleware(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Req)
  res := httptest.NewRecorder()
  var called bool
  next := func (c *core.Context, r http.ResponseWriter) {
    spec.Expect(c).ToEqual(context)
    spec.Expect(r).ToEqual(res)
    called = true
  }
  Run(context, res, next)
  spec.Expect(called).ToEqual(true)
}

func assertChunk(spec *gspec.S, chunk *core.Chunk, from int, to int) {
  spec.Expect(chunk.From).ToEqual(from)
  spec.Expect(chunk.To).ToEqual(to)
}
