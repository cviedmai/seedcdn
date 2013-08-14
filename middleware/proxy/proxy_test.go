package proxy

import (
  "testing"
  "seedcdn/core"
  "github.com/viki-org/gspec"
)

func TestCreatesARequestWithTheCorrectHostAndUrl(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Url("/test.json").Req)
  req := newRequest(context, &core.Config{Upstream: "s3.viki.com",})
  spec.Expect(req.Host).ToEqual("s3.viki.com")
  spec.Expect(req.URL.Path).ToEqual("/test.json")
}

func TestCreatesARequestWithTheCorrectRange(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Url("somefile.mp4").Req)
  context.Chunk = 3
  req := newRequest(context, &core.Config{RangedExtensions: map[string]bool{".mp4": true}})
  spec.Expect(req.Header.Get("range")).ToEqual("bytes=6291456-8388607")
}

func TestIgnoresTheRangeForNonRangeTypes(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Url("somefile.mp4").Req)
  context.Chunk = 3
  req := newRequest(context, new(core.Config))
  spec.Expect(req.Header.Get("range")).ToEqual("")
}
