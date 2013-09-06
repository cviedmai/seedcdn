package header

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestAnEmptyRange(t *testing.T) {
  spec := gspec.New(t)
  ranges := ParseRange("")
  spec.Expect(*ranges[0]).ToEqual(*&Range{0, 0})
}

func TestASimpleRange(t *testing.T) {
  spec := gspec.New(t)
  ranges := ParseRange("bytes=10-56")
  spec.Expect(len(ranges)).ToEqual(1)
  spec.Expect(*ranges[0]).ToEqual(*&Range{10, 56})
}

func TestMultipleRanges(t *testing.T) {
  spec := gspec.New(t)
  ranges := ParseRange("bytes=7-8,17-100")
  spec.Expect(len(ranges)).ToEqual(2)
  spec.Expect(*ranges[0]).ToEqual(*&Range{7, 8})
  spec.Expect(*ranges[1]).ToEqual(*&Range{17, 100})
}

func TestWithImplicitRangeStart(t *testing.T) {
  spec := gspec.New(t)
  ranges := ParseRange("bytes=-8")
  spec.Expect(len(ranges)).ToEqual(1)
  spec.Expect(*ranges[0]).ToEqual(*&Range{0, -8})
}

func TestWithImplicitRangeEnd(t *testing.T) {
  spec := gspec.New(t)
  ranges := ParseRange("bytes=193-")
  spec.Expect(len(ranges)).ToEqual(1)
  spec.Expect(*ranges[0]).ToEqual(*&Range{193, 0})
}

func TestComplexRange(t *testing.T) {
  spec := gspec.New(t)
  ranges := ParseRange("bytes=1-10,25-100,-10")
  spec.Expect(len(ranges)).ToEqual(3)
  spec.Expect(*ranges[0]).ToEqual(*&Range{1, 10})
  spec.Expect(*ranges[1]).ToEqual(*&Range{25, 100})
  spec.Expect(*ranges[2]).ToEqual(*&Range{0, -10})
}
