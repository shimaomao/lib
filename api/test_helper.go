package api

import (
	"fmt"
	// "github.com/speedland/lib/util"
	// "github.com/speedland/wcg/tools"
	// "io/ioutil"
	// "net/http"
	// "os"
	// "os/exec"
)

type TestApiServer struct {
	Port       int
	StorageDir string
}

func NewTestServer() *TestApiServer {
	var port int
	var storage string
	// TODO: Use mock server.
	server := &TestApiServer{
		Port:       port,
		StorageDir: storage,
	}
	return server
}

func (server *TestApiServer) Url(path string) string {
	return fmt.Sprintf("http://apps-dev.speedland.net/")
}

func (server *TestApiServer) Stop() {
	// defer os.RemoveAll(server.StorageDir)
	// defer server.Command.Process.Kill()
}
