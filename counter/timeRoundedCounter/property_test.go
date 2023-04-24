package timeRoundedCounter

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestPropertyCounter(t *testing.T) {
	now := time.Now()
	c := NewProperty("", time.Minute, 60)
	for i := -9; i <= 0; i++ {
		time.Sleep(time.Duration(rand.Int63n(int64(200 * time.Millisecond))))
		now = time.Now()

		ts := now.Add(time.Duration(i) * time.Minute)
		for j := 0; j < 10; j++ {
			property := fmt.Sprintf("test%d", 1+rand.Intn(5))
			c.Add(property, 1, ts)
		}
	}
	items := c.GetAll()
	assert.GreaterOrEqual(t, 10, len(items))

	printPropertyCountData(c.ItemProperties(), items)
}

func printPropertyCountData(properties []string, items []CounterItem) {
	util.Println("    " + strUtil.JoinStr(properties, "", func(i int, v string) string { return fmt.Sprintf("%-16s", v) }))

	for i, v := range items {
		arr := v.(PropertyCounterItem).Values()
		util.Println("[%d] %s", i, strUtil.Join(arr, "", func(i int) string { return fmt.Sprintf("%-16s", convertor.ToStringNoError(arr[i])) }))
	}
}
