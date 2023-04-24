package timeRoundedCounter

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestMRTCounter(t *testing.T) {
	now := time.Now()
	last := now
	c := NewStepMRT("", time.Minute, 60, 50*time.Millisecond, 100*time.Millisecond)
	for i := -9; i <= 0; i++ {
		time.Sleep(time.Duration(rand.Int63n(int64(200 * time.Millisecond))))
		now = time.Now()
		c.Add(1, now.Sub(last), now.Add(time.Duration(i)*time.Minute))
		last = now
	}
	items := c.GetAllMRT()
	assert.GreaterOrEqual(t, 10, len(items))
	printMRTCountData(c.ItemProperties(), items)
}

func TestMRTCounter_Simple(t *testing.T) {
	now := time.Now()
	last := now
	c := NewStepMRT("", time.Minute, 60)
	for i := -9; i <= 0; i++ {
		time.Sleep(time.Duration(rand.Int63n(int64(100 * time.Millisecond))))
		now = time.Now()
		c.Add(1, now.Sub(last), now.Add(time.Duration(i)*time.Minute))
		last = now
	}
	items := c.GetAllMRT()
	assert.GreaterOrEqual(t, 10, len(items))
	printMRTCountData(c.ItemProperties(), items)
}

func printMRTCountData(properties []string, items []MRTCounterItem) {
	if len(items) == 0 {
		return
	}

	util.Println("    " + strUtil.JoinStr(properties, "", func(i int, v string) string { return fmt.Sprintf("%-16s", v) }))

	for i, v := range items {
		arr := v.(MRTCounterItem).MRTValues()
		util.Println("[%d] %s", i, strUtil.Join(arr, "", func(i int) string { return fmt.Sprintf("%-16s", fmt.Sprintf("%d, %s", arr[i].Count, arr[i].MRT)) }))
	}
}
