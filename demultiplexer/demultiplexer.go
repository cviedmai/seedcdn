package demultiplexer

import (
  "sync"
  "seedcdn/core"
  "seedcdn/middleware/proxy"
)

var (
  lock sync.RWMutex
  masters = make(map[string] *Master)
)

type Handler func(payload *Payload)

func Demultiplex(context *core.Context, slaveHandler Handler, masterHandler Handler) {
  master, new := getMaster(context.Key())
  if new == true {
    go master.Run(proxy.Run(context))
    go func() {
      c := make(chan *Payload, 1)
      master.Observed(c)
      for {
        payload := <- c
        masterHandler(payload)
        if payload.Finished { return }
      }
    }()
  }
  c := make(chan *Payload, 1)
  master.Observed(c)
  for {
    payload := <- c
    slaveHandler(payload)
    if payload.Finished { return }
  }
}

func Cleanup(key string) {
  lock.Lock()
  defer lock.Unlock()
  delete(masters, key)
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
