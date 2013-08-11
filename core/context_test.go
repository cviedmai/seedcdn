package core

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestDefaultsToTheFirstChunk(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Req)
  spec.Expect(c.Chunk).ToEqual(0)
}

func TestSetsTheRightChuckWhenSmallRange(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Header("range", "bytes=7000000-7000001").Req)
  spec.Expect(c.Chunk).ToEqual(3)
}

func TestSetsTheRightChuckWhenLargeRange(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Header("range", "bytes=3000000-9000000").Req)
  spec.Expect(c.Chunk).ToEqual(1)
}

func TestTheKeyIncludesAllNecessaryParts(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Header("range", "bytes=30000000-90000000").Url("/test.json?page=1").Req)
  spec.Expect(c.Key()).ToEqual("/test.json?page=1_14")
}
