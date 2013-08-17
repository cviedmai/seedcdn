package chunker

import (
  "fmt"
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
  spec.Expect(context.Chunks[0]).ToEqual(buildChunk(0, 0, "", ""))
}

func TestHandlesARangeWithNoLimit(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-").Req)
  Run(context, nil, core.NoopMiddleware)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  spec.Expect(context.Chunks[0]).ToEqual(buildChunk(10, 0, "", ""))
}

func TestHandlesARangeWithinASingleChunk(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-2000000").Req)
  Run(context, nil, core.NoopMiddleware)
  fmt.Println(context.Chunks)
  spec.Expect(len(context.Chunks)).ToEqual(1)
  spec.Expect(context.Chunks[0]).ToEqual(buildChunk(10, 2000000, "6666cd76f96956469e7be39d750cc7d9_0", "/tmp/66/6666/6666cd76f96956469e7be39d750cc7d9/6666cd76f96956469e7be39d750cc7d9_0"))
}

func TestHandlesARangeAcrossMultipleChunks(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Header("Range", "bytes 10-12000000").Req)
  Run(context, nil, core.NoopMiddleware)
  fmt.Println(context.Chunks)
  spec.Expect(len(context.Chunks)).ToEqual(3)
  spec.Expect(context.Chunks[0]).ToEqual(buildChunk(10, 5242879, "6666cd76f96956469e7be39d750cc7d9_0", "/tmp/66/6666/6666cd76f96956469e7be39d750cc7d9/6666cd76f96956469e7be39d750cc7d9_0"))
  spec.Expect(context.Chunks[1]).ToEqual(buildChunk(5242880, 10485759, "6666cd76f96956469e7be39d750cc7d9_1", "/tmp/66/6666/6666cd76f96956469e7be39d750cc7d9/6666cd76f96956469e7be39d750cc7d9_1"))
  spec.Expect(context.Chunks[2]).ToEqual(buildChunk(10485760, 12000000, "6666cd76f96956469e7be39d750cc7d9_2", "/tmp/66/6666/6666cd76f96956469e7be39d750cc7d9/6666cd76f96956469e7be39d750cc7d9_2"))
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

func buildChunk(from int, to int, key string, dataFile string) core.Chunk {
  return *&core.Chunk{
    From: from,
    To: to,
    Key: key,
    DataFile: dataFile,
  }
}
