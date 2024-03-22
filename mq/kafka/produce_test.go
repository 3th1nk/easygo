package kafka

import (
	"fmt"
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestKafka_Produce(t *testing.T) {
	kafka, err := New(strings.Split(_test.KafkaAddr, ","))
	runtimeUtil.PanicIfError(err)

	const loop = 100
	for i := 0; i < loop; i++ {
		partition, offset, err := kafka.Produce("test", fmt.Sprintf("test-%d", i))
		assert.NoError(t, err)
		t.Logf("%v: %v", partition, offset)
	}
}

func TestKafka_Produce_Perf(t *testing.T) {
	kafka, err := New(strings.Split(_test.KafkaAddr, ","))
	runtimeUtil.PanicIfError(err)

	_test.PerfIf(func(i int) bool {
		_, _, err = kafka.Produce("test", fmt.Sprintf("test-%d", i))
		return err == nil
	}, _test.PerfOptions{Dur: 3 * time.Second, Name: "Sync-Goroutine-1", Goroutine: 1})

	_test.PerfIf(func(i int) bool {
		_, _, err = kafka.Produce("test", fmt.Sprintf("test-%d", i))
		return err == nil
	}, _test.PerfOptions{Dur: 3 * time.Second, Name: "Sync-Goroutine-N", Goroutine: 100})
}
