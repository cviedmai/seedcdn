package seedcdn

import (
  "net/http"
  "seedcdn/header"
)

type Handler struct{}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  header.ParseRange(req.Header.Get("range"))
}
