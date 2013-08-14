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

func TestTheContextsDirectories(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Url("/test.json").Req)
  spec.Expect(c.Dir).ToEqual("/tmp/01/0196/0196f4b7a30827487e3272e9499749e9/")
  spec.Expect(c.TempDir).ToEqual("/tmp/tmp/")
}

func TestTheContextsFiles(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Url("/over9000.json").Req)
  spec.Expect(c.HeaderFile).ToEqual("/tmp/d2/d21f/d21fd0eba7a9f34038292d38c8ff4837/d21fd0eba7a9f34038292d38c8ff4837_0.hdr")
  spec.Expect(c.DataFile).ToEqual("/tmp/d2/d21f/d21fd0eba7a9f34038292d38c8ff4837/d21fd0eba7a9f34038292d38c8ff4837_0.dat")
}

func TestTheContextsKey(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Header("range", "bytes=30000000-90000000").Url("/test.json?page=1").Req)
  spec.Expect(c.Key).ToEqual("0196f4b7a30827487e3272e9499749e9_14")
}
