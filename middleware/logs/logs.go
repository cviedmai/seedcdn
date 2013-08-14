package logs

import (
  "log"
  "time"
  "net/http"
  "seedcdn/core"
  "github.com/garyburd/redigo/redis"
)

var pool = &redis.Pool {
  MaxIdle: 5,
  IdleTimeout: time.Minute * 5,
  Dial: func () (redis.Conn, error) {
    config := core.GetConfig()
    conn, err := redis.Dial(config.RedisProtocol, config.Redis)
    if err != nil  { return nil, err }
    if _, err = conn.Do("select", config.RedisDB); err != nil { return nil, err }
    return conn, nil
  },
}

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  go run(context)
  next(context, res)
}

func run(context *core.Context) {
  path := context.Req.URL.Path
  conn := pool.Get()
  defer conn.Close()
  conn.Send("multi")
  conn.Send("zincrby", "hits", 1, path)
  conn.Send("zadd", "last", time.Now().Unix(), path)
  if _, err := conn.Do("exec"); err != nil {
    log.Println("log: ", err)
  }
}
