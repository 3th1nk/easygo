package influxdb

import (
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Query(t *testing.T) {
	c := Between("time", "2024-05-16 18:40:00", "2024-05-16 18:46:00")
	t.Log(c.String())

	q := influx.NewQuery()
	q.Select("*").From(testDB, testRP.Name, "measurement-1").Where(c)
	t.Log(q.String())

	res, err := q.Do()
	assert.NoError(t, err)
	t.Log(jsonUtil.MustMarshalToStringIndent(res))
}
