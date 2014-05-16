package api

import (
	"encoding/json"
	"fmt"
	"github.com/speedland/lib"
	"github.com/speedland/wcg"
	"io/ioutil"
	"net/http"
	"reflect"
)

type ErrUnexpectedStatusCode struct {
	Expect int
	Actual int
}

func (e *ErrUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("Get the unexpected status code %d (expected: %d)", e.Actual, e.Expect)
}

type ErrDecodingJson struct {
	Body    string
	Type    reflect.Type
	Message string
}

func (e *ErrDecodingJson) Error() string {
	return fmt.Sprintf("Could not decode JSON string as %s: %s (%q)", e.Type, e.Message, e.Body)
}

type ApiClient struct {
	http.Client
}

type apiRoundTripper struct {
	http.RoundTripper
}

func (rt *apiRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	logger := wcg.NewLogger(nil)
	req.Header.Set("X-SPEEDLAND-API-TOKEN", lib.Config.Token)
	logger.Debug("[Api] %s %s", req.Method, req.URL.Path)
	return rt.RoundTripper.RoundTrip(req)
}

func newApiRoundTripper() http.RoundTripper {
	rt := new(apiRoundTripper)
	rt.RoundTripper = http.DefaultTransport
	return rt
}

var DefaultApiTransport = newApiRoundTripper()
var DefaultApiClient = &ApiClient{
	http.Client{
		Transport: DefaultApiTransport,
	},
}

func (c *ApiClient) Ping() (map[string]interface{}, error) {
	endpoint := buildUrl("/api/auth/me")
	resp, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}
	var me map[string]interface{}
	err = handleAsJson(resp, &me)
	if err != nil {
		return nil, err
	} else {
		return me, nil
	}
}

func Ping() (map[string]interface{}, error) {
	return DefaultApiClient.Ping()
}

func handleAsJson(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, v)
	if err != nil {
		return &ErrDecodingJson{
			Body:    string(buff),
			Type:    reflect.TypeOf(v),
			Message: err.Error(),
		}
	}
	return nil
}

func checkStatusCode(expect, actual int) error {
	if expect != actual {
		return &ErrUnexpectedStatusCode{
			Expect: expect,
			Actual: actual,
		}
	}
	return nil
}

func buildUrl(path string) string {
	return lib.Config.Endpoint.String() + path
}
