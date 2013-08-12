package demultiplexer

import (
  "sync"
  "net/http"
  "seedcdn/core"
)

var (
  lock sync.RWMutex
  masters = make(map[string] *Master)
)

type Proxy func (context *core.Context) http.Response

func Demultiplex(context *core.Context, proxy Proxy) {
  master, new := getMaster(context.Key())
  if new == true {
    master.Run(proxy(context))
    return
  }
  c := make(chan []byte, IDEAL_CHUNK_COUNT)
  master.observers <- c
  sync := <- master.sync
  println(sync.Status())
  //res header
  for {
    chunklet := <- c
    if len(chunklet) == 0 { break }
    //res.write(chunklet)
  }
}

func getMaster(key string) (*Master, bool) {
  lock.RLock()
  master, ok := masters[key]
  lock.RUnlock()
  if ok == true { return master, false }

  lock.Lock()
  defer lock.Unlock()
  master, ok = masters[key]
  if ok == true { return master, false }

  master = New(key)
  masters[key] = master
  return master, true
}
