package timeRoundedCounter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestPropertyMRTCounter(t *testing.T) {
	now := time.Now()
	last := now
	c := NewPropertyMRT("", time.Minute, 60)
	for i := -9; i <= 0; i++ {
		time.Sleep(time.Duration(rand.Int63n(int64(200 * time.Millisecond))))
		now = time.Now()

		d, ts := now.Sub(last), now.Add(time.Duration(i)*time.Minute)
		for j := 0; j < 5; j++ {
			property := fmt.Sprintf("test%d", 1+rand.Intn(8))
			c.Add(property, 1, d, ts)
		}

		last = now
	}
	items := c.GetAllMRT()
	assert.GreaterOrEqual(t, 10, len(items))
	printMRTCountData(c.ItemProperties(), items)
}
