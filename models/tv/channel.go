package tv

import (
	"fmt"
	v "github.com/speedland/wcg/validation"
)

type TvChannel struct {
	Cid  string `json:"cid"`
	Sid  string `json:"sid"`
	Name string `json:"name"`
}

var TvChannelValidator = v.NewObjectValidator()

func init() {
	TvChannelValidator.Field("Sid").Required()
	TvChannelValidator.Field("Cid").Required()
	TvChannelValidator.Field("Name").Required()
}

func (c *TvChannel) Key() string {
	return fmt.Sprintf("%s.%s", c.Cid, c.Sid)
}

func (c *TvChannel) String() string {
	return fmt.Sprintf("<TvChannelConfig : [%s] %s>", c.Key(), c.Name)
}
