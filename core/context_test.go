package core

import (
  "testing"
  "github.com/viki-org/gspec"
)

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
}