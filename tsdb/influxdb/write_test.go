package influxdb

import (
	"fmt"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func makePoints(n int) []*Point {
	points := make([]*Point, 0, n)
	for i := 0; i < n; i++ {
		points = append(points, &Point{
			Measurement: fmt.Sprintf("measurement-%d", i%10),
			Tags: map[string]interface{}{
				"tag1": fmt.Sprintf("tag1-%d", i%5),
				"tag2": fmt.Sprintf("tag2-%d", i%5),
				"tag3": fmt.Sprintf("tag3-%d", i%5),
			},
			Values: map[string]interface{}{
				"value": i % 100,
			},
		})
	}

	return points
}

func TestClient_WriteAndQuery(t *testing.T) {
	err := initTestDbRp()
	assert.NoError(t, err)

	points := makePoints(10)
	err = influx.Write(testDB, testRP.Name, points, true)
	assert.NoError(t, err)

	cond := Or(Expr("tag1", "=", "tag1-1"), Expr("tag1", "=", "tag1-2"))
	q := influx.NewQuery().From(testDB, testRP.Name, "measurement-1", "measurement-2").Where(cond).Limit(10)
	t.Log(q.String())

	res, err := q.Do()
	assert.NoError(t, err)
	t.Log(jsonUtil.MustMarshalToStringIndent(res))
}

func TestClient_WriteReliability(t *testing.T) {
	err := initTestDbRp()
	assert.NoError(t, err)

	st := time.Now()
	defer func() {
		t.Logf("time cost: %v", time.Since(st))
	}()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			points := makePoints(1000)
			_ = influx.Write(testDB, testRP.Name, points, false)
		}()
	}
	wg.Wait()
	influx.Close()

	influx.showWriteCount()
}

func BenchmarkClient_WritePerf(b *testing.B) {
	err := initTestDbRp()
	assert.NoError(b, err)

	points := makePoints(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		influx.startAsyncWrite()
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = influx.Write(testDB, testRP.Name, points, false)
			}()
		}
		wg.Wait()
		influx.stopAsyncWrite()
	}

}
