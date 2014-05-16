// SPEEDLAND Common Library
package lib

import (
	"github.com/speedland/wcg"
	"net/url"
)

type config struct {
	Endpoint *url.URL `ini:"endpoint" default:"http://apps-dev.speedland.net/"`
	Token    string   `ini:"token" default:""`
}

var Config = &config{}

func init() {
	wcg.RegisterProcessConfig(Config, "speedland", nil)
}
