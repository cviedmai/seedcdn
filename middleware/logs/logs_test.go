package logs

import (
  "time"
  "testing"
  "net/http"
  "seedcdn/core"
  "net/http/httptest"
  "github.com/viki-org/gspec"
  "github.com/garyburd/redigo/redis"
)

func TestCallsTheNextMiddleware(t *testing.T) {
  spec := gspec.New(t)
  context := core.NewContext(gspec.Request().Req)
  res := httptest.NewRecorder()
  var called bool
  next := func (c *core.Context, r http.ResponseWriter) {
    spec.Expect(c).ToEqual(context)
    spec.Expect(r).ToEqual(res)
    called = true
  }
  Run(context, res, next)
  spec.Expect(called).ToEqual(true)
}

func TestLogsStatisticsToRedis(t *testing.T) {
  conn := core.RedisPool.Get()
  defer cleanup(conn)
  spec := gspec.New(t)
  Run(core.NewContext(gspec.Request().Url("/something/funny.txt").Req), nil, core.NoopMiddleware)
  time.Sleep(time.Second * 1)
  spec.Expect(redis.Int(conn.Do("zscore", "hits", "/something/funny.txt"))).ToEqual(1)
}

func cleanup(conn redis.Conn) {
  conn.Do("flushdb")
  conn.Close()
}
