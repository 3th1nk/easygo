package rabbitMQ

import (
	"fmt"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProducer_Pool(t *testing.T) {
	const queueName = "test_3_1"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	testWorkerLoops(1000, time.Minute, func(worker, loop int) {
		err := queue.Produce("", queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})
}

func TestProducer_Queue(t *testing.T) {
	const queueName = "test_4_1"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	producer, err := queue.NewProducer()
	assert.NoError(t, err)

	testWorkerLoops(1000, time.Minute, func(worker, loop int) {
		err := producer.SendToQueue(queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})

}

func TestProducer_GroupQueue(t *testing.T) {
	const queueName = "test_5_1"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	producer, err := queue.NewGroupProducer(10)
	assert.NoError(t, err)

	testWorkerLoops(1000, time.Minute, func(worker, loop int) {
		err := producer.SendToQueue(queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})

}
