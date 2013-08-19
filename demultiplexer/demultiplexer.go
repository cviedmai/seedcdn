package demultiplexer

import (
  "sync"
  "seedcdn/core"
)

var (
  lock sync.RWMutex
  masters = make(map[string] *Master)
)

type Handler func(payload *Payload)

func Demultiplex(context *core.Context, chunk *core.Chunk, slaveHandler Handler, masterHandler Handler) int {
  master, new := getMaster(context.Key)
  if new == true {
    res, err := download(context, chunk)
    go master.Run(res, err, masterHandler)
  }
  c := make(chan *Payload)
  master.Observed(c)
  for {
    payload := <- c
    slaveHandler(payload)
    if payload.Finished { return payload.ContentLength }
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

  master = &Master{key: key,}
  masters[key] = master
  return master, true
}
