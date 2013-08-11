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

func Demultiplex(context *core.Context) {
  master, new := getMaster(context.Key())
  if new == true {
    master.Run(proxy(context))
    return
  }
  //todo attach to master
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

  master = &Master{
    key: key,
    header: make(http.Header, len(proxyHeaders)),
  }
  masters[key] = master
  return master, true
}
