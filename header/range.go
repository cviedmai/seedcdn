package header

import (
  "strconv"
)

type Range struct {
  From int
  To int
}

var fullRange = []Range{*&Range{0, 0}}

func ParseRange(raw string) []Range {
  if len(raw) == 0 { return fullRange }
  var ranges []Range
  length := len(raw)
  start := 6 //bytes=
  split := 0
  for i := start; i < length; i++ {
    c := raw[i]
    if c == ',' {
      ranges = append(ranges, createRange(raw[start:i], split - start))
      start = i + 1
    }
    if c == '-' {
      split = i
    }
  }
  return append(ranges, createRange(raw[start:], split - start))
}

func createRange(raw string, split int) Range {
  var from, to int
  var relative bool
  if split == 0 {
    from = 0
    relative = true
  } else {
    from, _ = strconv.Atoi(raw[:split])
  }
  if split == len(raw) {
    to = 0
  } else {
    to, _ = strconv.Atoi(raw[split+1:])
    if relative { to *= -1 }
  }
  return Range{from, to}
}
