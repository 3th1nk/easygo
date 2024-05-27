package influxdb

import (
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util/logs"
)

var (
	testDB = "test"
	testRP = &RetentionPolicy{
		Name:        "test_rp",
		Duration:    "7d",
		Replication: 1,
	}
)

var (
	influx *Client
)

func init() {
	influx = NewClient(_test.InfluxAddr,
		WithLogger(logs.Stdout(logs.LevelAll)),
		withDebugger(true),
	)
}
