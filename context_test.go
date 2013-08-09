package seedcdn

import (
  "testing"
  "seedcdn/gspec"
)

func TestDefaultsToTheFirstChunk(t *testing.T) {
  spec := gspec.New(t)
  c := NewContext(gspec.Request().Req())
  spec.Expect(c.chunk).ToEqual(0)
}
