package api

import (
	"github.com/speedland/lib"
	"github.com/speedland/wcg"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	assert := wcg.NewAssert(t)
	server := NewTestServer()
	defer server.Stop()
	lib.Config.Endpoint, _ = url.Parse(server.Url(""))
	lib.Config.Token = "dummy"
	_, err := Ping()
	assert.Nil(err, "Ping should not return an erorr.")

	lib.Config.Endpoint, _ = url.Parse("http://localhost:3001")
	lib.Config.Token = "dummy"
	_, err = Ping()
	assert.NotNil(err, "Ping should return an error when it could not reach the server.")

}
