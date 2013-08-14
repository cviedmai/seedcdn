package core

import (
  "os"
  "fmt"
  "time"
  "runtime"
  "sync/atomic"
)

type stats struct {
  requests uint64
  cacheHits uint64
  cacheMisses uint64
}

var Stats = new(stats)

func init() { go snapshot() }

func (s *stats) Request() { atomic.AddUint64(&s.requests, 1) }
func (s *stats) CacheHit() { atomic.AddUint64(&s.cacheHits, 1) }
func (s *stats) CacheMiss() { atomic.AddUint64(&s.cacheMisses, 1) }

func snapshot() {
  var last = new(stats)
  for {
    time.Sleep(time.Minute)
    requests := atomic.LoadUint64(&Stats.requests)
    cacheHits := atomic.LoadUint64(&Stats.cacheHits)
    cacheMisses := atomic.LoadUint64(&Stats.cacheMisses)

    ch := cacheHits - last.cacheHits
    ct := cacheMisses - last.cacheMisses + ch
    if ct == 0 { ct = 1 }

    s := []byte(fmt.Sprintf(`{"requests":%d,"cacheRatio":%d,"goroutines":%d}`, requests - last.requests, ch * 100 / ct, runtime.NumGoroutine()))

    last.requests = requests
    last.cacheHits = cacheHits
    last.cacheMisses = cacheMisses

    file, _ := os.Create("stats.json")
    file.Write(s)
    file.Close()
  }
}
