package purge

import (
  "os"
  "strings"
  "net/http"
  "seedcdn/core"
)

func Run (context *core.Context, res http.ResponseWriter, next core.Middleware) {
  if context.Req.Method != "PURGE" {
    next(context, res)
    return
  }

  ip := context.Req.RemoteAddr[0:strings.LastIndex(context.Req.RemoteAddr, ":")]
  if core.GetConfig().PurgeWhiteList[ip] == false {
    res.WriteHeader(401)
    return
  }

  if err := os.RemoveAll(context.Dir); err != nil {
    res.Write([]byte(err.Error()))
    res.WriteHeader(500)
  } else  {
    res.WriteHeader(200)
  }
}
