package proxy

import (
  "net"
  "path"
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

func Run(context *core.Context) (*http.Response, error) {
  request := newRequest(context, core.GetConfig())
  return transport.RoundTrip(request)
}

func newRequest(context *core.Context, config *core.Config) *http.Request {
  header := make(http.Header)
  if config.RangedExtensions[path.Ext(context.Req.URL.Path)] == true {
    from := context.Chunk * int(core.CHUNK_SIZE)
    to := from + int(core.CHUNK_SIZE) - 1
    header.Set("Range", "bytes=" + strconv.Itoa(from) + "-" + strconv.Itoa(to))
  }
  u := context.Req.URL
  return &http.Request{
    Close: false,
    Host: config.Upstream,
    Method: "GET",
    Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
    Header: header,
    URL: &url.URL{
      Scheme: "http",
      Opaque: u.Opaque,
      Host: config.Upstream,
      Path: config.Prefix + u.Path,
      RawQuery: u.RawQuery,
    },
  }
}

func dial(network string, address string) (net.Conn, error) {
  separator := strings.LastIndex(address, ":")
  ip, _ := dns.FetchOneString(address[:separator])
  return net.Dial("tcp", ip + address[separator:])
}
