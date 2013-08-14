package purge

import (
  "os"
  "testing"
  "net/http"
  "seedcdn/core"
  "net/http/httptest"
  "github.com/viki-org/gspec"
)

func TestGoesToTheNextMiddlewareForNonPurgeRequests(t *testing.T) {
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

func TestDeletesTheFilesDirectory(t *testing.T) {
  spec := gspec.New(t)
  dir := setupTempFolder(spec)
  defer cleanupTempFolder()

  res := httptest.NewRecorder()
  Run(core.NewContext(gspec.Request().Method("purge").Req), res, nil)
  spec.Expect(res.Code).ToEqual(200)
  _, err := os.Stat(dir)
  spec.Expect(err).ToNotBeNil()
}

func TestOnlyAllowsWhitelistedIpsToPurge(t *testing.T) {
  spec := gspec.New(t)
  dir := setupTempFolder(spec)
  defer cleanupTempFolder()

  res := httptest.NewRecorder()
  Run(core.NewContext(gspec.Request().Method("purge").RemoteAddr("23.33.24.55:4343").Req), res, nil)
  spec.Expect(res.Code).ToEqual(401)
  spec.Expect(res.Body.Len()).ToEqual(0)
  _, err := os.Stat(dir)
  spec.Expect(err).ToBeNil()
}

func setupTempFolder(spec *gspec.S) string {
  dir := "/tmp/66/6666/6666cd76f96956469e7be39d750cc7d9/"
  spec.Expect(os.MkdirAll(dir, 0744)).ToBeNil()
  file, _ := os.Create(dir + "test")
  file.Close()
  return dir
}

func cleanupTempFolder() {
  os.RemoveAll("/tmp/66/")
}
