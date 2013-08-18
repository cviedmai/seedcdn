package demultiplexer

import (
  "net"
  "time"
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

func download(context *core.Context, chunk *core.Chunk) (*http.Response, error) {
  request := newRequest(context, chunk, core.GetConfig())
  return transport.RoundTrip(request)
}

func newRequest(context *core.Context, chunk *core.Chunk, config *core.Config) *http.Request {
  u := context.Req.URL
  return &http.Request{
    Close: false,
    Host: config.Upstream,
    Method: "GET",
    Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
    Header: chunk.Header,
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
