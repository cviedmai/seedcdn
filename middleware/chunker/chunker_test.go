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
  spec.Expect(len(context.Chunks)).ToEqual(1)
  spec.Expect(*context.Chunks[0]).ToEqual(buildChunk(0, 0, 0))
}

func TestHandlesARangeWithNoLimit(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  spec.Expect(*context.Chunks[0]).ToEqual(buildChunk(0, 0, 0))
}

func TestHandlesARangeWithNoLimitAfterTheFirstChunk(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 7000000-").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  spec.Expect(*context.Chunks[0]).ToEqual(buildChunk(5242880, 0, 1))
}

func TestHandlesARangeWithinASingleChunk(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-2000000").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  spec.Expect(*context.Chunks[0]).ToEqual(buildChunk(0, 5242879, 0))
}

func TestHandlesARangeAcrossMultipleChunks(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-12000000").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(len(context.Chunks)).ToEqual(3)
  spec.Expect(*context.Chunks[0]).ToEqual(buildChunk(0, 5242879, 0))
  spec.Expect(*context.Chunks[1]).ToEqual(buildChunk(5242880, 10485759, 1))
  spec.Expect(*context.Chunks[2]).ToEqual(buildChunk(10485760, 15728639, 2))
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

func buildChunk(from int, to int, n int) core.Chunk {
  return *&core.Chunk{
    From: from,
    To: to,
    N: n,
  }
}
