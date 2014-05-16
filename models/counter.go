package models

import (
	"fmt"
	influx "github.com/influxdb/influxdb-go"
	"github.com/speedland/wcg"
	"net/http"
	"os"
)

type influxConfig struct {
	Host     string `ini:"host" default:"sandbox.influxdb.com"`
	Port     int    `ini:"port" default:"8086"`
	Username string `ini:"username" default:"yssk22"`
	Password string `ini:"password" default:"passw0rd"`
	Database string `ini:"database" default:"yssk22-test"`
	IsSecure bool   `ini:"is_secure" default:"false"`
}

var InfluxConfig = new(influxConfig)
var influxSource string

type CounterClient struct {
	influxClient *influx.Client
}

func NewCounterClient(httpClient *http.Client) (*CounterClient, error) {
	client, err := influx.NewClient(&influx.ClientConfig{
		Host:       fmt.Sprintf("%s:%d", InfluxConfig.Host, InfluxConfig.Port),
		Username:   InfluxConfig.Username,
		Password:   InfluxConfig.Password,
		Database:   InfluxConfig.Database,
		IsSecure:   InfluxConfig.IsSecure,
		HttpClient: httpClient,
	})
	return &CounterClient{influxClient: client}, err
}

func (cc *CounterClient) Post(name string, value interface{}) error {
	series := &influx.Series{
		Name:    name,
		Columns: []string{"source", "value"},
		Points: [][]interface{}{
			[]interface{}{influxSource, value},
		},
	}
	return cc.influxClient.WriteSeries([]*influx.Series{series})
}

func (cc *CounterClient) PostMany(name string, table map[string]interface{}) error {
	cols := []string{"source"}
	vals := []interface{}{influxSource}
	for k, v := range table {
		cols = append(cols, k)
		vals = append(vals, v)
	}
	series := &influx.Series{
		Name:    name,
		Columns: cols,
		Points: [][]interface{}{
			vals,
		},
	}
	return cc.influxClient.WriteSeries([]*influx.Series{series})
}

func init() {
	wcg.RegisterProcessConfig(InfluxConfig, "speedland.models.influxconfig", nil)
	influxSource, _ = os.Hostname()
}
