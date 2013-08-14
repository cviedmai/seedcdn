package core

import (
  "testing"
  "github.com/viki-org/gspec"
)

func TestLoadsAConfiguration(t *testing.T) {
  spec := gspec.New(t)
  loadConfig([]byte(`{"listen":"1.123.58.13:9001", "upstream":"its.over.net"}`))
  config := GetConfig()
  spec.Expect(config.Listen).ToEqual("1.123.58.13:9001")
  spec.Expect(config.Upstream).ToEqual("its.over.net")
}
