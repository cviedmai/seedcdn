package proxy

import (
  "io"
  "log"
  "net"
  "time"
  "strconv"
  "strings"
  "net/url"
  "net/http"
  "seedcdn/core"
  "github.com/viki-org/dnscache"
)

var dns = dnscache.New(time.Minute * 2)
var transport = &http.Transport{
  MaxIdleConnsPerHost: 32,
  DisableKeepAlives: false,
  Dial: dial,
}

func Run(context *core.Context, res http.ResponseWriter, next core.Middleware) {
  request := newRequest(context, core.GetConfig())
  r, err := transport.RoundTrip(request)
  if r != nil && r.Body != nil { defer r.Body.Close() }
  if err != nil { log.Println("upstream error: ", err) }
  io.Copy(res, r.Body)
}

func newRequest(context *core.Context, config *core.Config) *http.Request {
  from := context.Chunk * int(core.CHUNKSIZE)
  to := from + int(core.CHUNKSIZE) - 1

  u := context.Req.URL
  return &http.Request{
    Close: false,
    Host: config.Upstream,
    Method: "GET",
    Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
    Header: http.Header{"Range": []string{"bytes=" + strconv.Itoa(from) + "-" + strconv.Itoa(to)}},
    URL: &url.URL{
      Scheme: "http",
      Opaque: u.Opaque,
      Host: config.Upstream,
      Path: u.Path,
      RawQuery: u.RawQuery,
    },
  }
}

func dial(network string, address string) (net.Conn, error) {
  separator := strings.LastIndex(address, ":")
  ip, _ := dns.FetchOneString(address[:separator])
  return net.Dial("tcp", ip + address[separator:])
}
