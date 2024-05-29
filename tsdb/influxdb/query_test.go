package influxdb

import (
	"github.com/3th1nk/easygo/util/jsonUtil"
	"testing"
)

func TestClient_Query(t *testing.T) {
	q := influx.NewQuery()
	cond := NewCond()
	q.Select("*").From(testDB, "", "measurement1").Where(cond)
	t.Log(q.String())

	res, err := q.Do()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(jsonUtil.MustMarshalToStringIndent(res))
}
