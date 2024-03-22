package rabbitMQ

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/logs"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/3th1nk/easygo/util/timeUtil"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_Produce_1(t *testing.T) {
	const queueName = "test_1"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	testWorkerLoops(1, 3*time.Second, func(worker, loop int) {
		err := queue.Produce("", queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})
}

func Test_Produce_10(t *testing.T) {
	const queueName = "test_2"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	testWorkerLoops(10, 3*time.Second, func(worker, loop int) {
		err := queue.Produce("", queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})
}

func Test_Produce_100(t *testing.T) {
	const queueName = "test_3"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	testWorkerLoops(100, 3*time.Second, func(worker, loop int) {
		err := queue.Produce("", queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})
}

func Test_Produce_1000(t *testing.T) {
	const queueName = "test_4"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	testWorkerLoops(1000, 15*time.Second, func(worker, loop int) {
		err := queue.Produce("", queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
	})
}

func Test_Produce_1000_1(t *testing.T) {
	const queueName = "test_4_1"

	declareQueue(t, queueName)
	defer removeQueue(t, queueName)

	queue := newTestQueue(runtimeUtil.CallerFunc(0))
	defer queue.Stop(0)

	start := time.Now()
	err := queue.Do(func(ch *amqp.Channel) error {
		a := time.Now()
		testWorkerLoops(1000, 15*time.Second, func(worker, loop int) {
			err := ch.Publish("", queueName, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         []byte(fmt.Sprintf("%d-%d", worker, loop)),
			})
			assert.NoError(t, err)
		})

		fmt.Println("---->", time.Since(a).String())
		return nil
	})
	fmt.Println("----", time.Since(start).String())
	assert.NoError(t, err)
}

func TestConcurrence(t *testing.T) {
	const queueCount, produceWorker, consumeWorker, consumeConcurrent = 10, 100, 100, 1000
	pool := channelPool
	queueOption := Options{NoAutoDeclare: true, NoProduceLog: true, NoConsumeLog: true, Logger: logs.Empty}

	// queue
	queueNameList := make([]string, queueCount)
	for i := 0; i < queueCount; i++ {
		queueNameList[i] = fmt.Sprintf("test_%d_%d", queueCount, i)
	}
	declareQueue(t, queueNameList...)
	defer removeQueue(t, queueNameList...)

	// producer
	producerList := make([]*Queue, produceWorker)
	for i := 0; i < produceWorker; i++ {
		queue := newTestQueue(fmt.Sprintf("producer-%d", i), pool, queueOption)
		queue.SetLogger(logs.Empty)
		producerList[i] = queue
	}
	defer func() {
		for _, queue := range producerList {
			queue.Stop(0)
		}
	}()

	// consumer
	consumerList := make([]*Queue, consumeWorker)
	for i := 0; i < consumeWorker; i++ {
		queue := newTestQueue(fmt.Sprintf("consumer-%d", i), pool, queueOption)
		queue.SetLogger(logs.Empty)
		consumerList[i] = queue
	}
	defer func() {
		for _, queue := range consumerList {
			queue.Stop(0)
		}
	}()

	produce, consume, lastProduce, lastConsume, lastTS := int32(0), int32(0), int32(0), int32(0), time.Now()
	ticker := timeUtil.NewTicker(500*time.Millisecond, 500*time.Millisecond, func(now time.Time) {
		sec := now.Sub(lastTS).Seconds()
		produceTps, consumeTps := float64(produce-lastProduce)/sec, float64(consume-lastConsume)/sec
		util.Println("    [%s] produce=%d, consume=%d, produceTps=%.0f, consumeTps=%.0f", now.Format("15:04:05.000"), produce, consume, produceTps, consumeTps)
		lastProduce, lastConsume, lastTS = produce, consume, now
	})
	defer func() {
		ticker.Trigger()
		ticker.Stop(0)
	}()

	for i, queue := range consumerList {
		queueName := queueNameList[i%queueCount]
		_, err := queue.Consume(queueName, consumeConcurrent, func(data string) (ack bool) {
			atomic.AddInt32(&consume, 1)
			return true
		})
		assert.NoError(t, err)
	}
	testWorkerLoops(produceWorker, 10*time.Second, func(worker, loop int) {
		queue, queueName := producerList[worker], queueNameList[worker%queueCount]
		err := queue.Produce("", queueName, fmt.Sprintf("%d-%d", worker, loop))
		assert.NoError(t, err)
		atomic.AddInt32(&produce, 1)
	}, true)
	util.Println("    [%s] -------------------------- PRODUCE END --------------------------", time.Now().Format("15:04:05.000"))

	for consume < produce {
		time.Sleep(100 * time.Millisecond)
	}
	ticker.Trigger()
}

func testWorkerLoops(worker int, dur time.Duration, f func(worker, loop int), noDebug ...bool) {
	loop, done := int32(0), int32(0)

	start := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(worker)
	for i := 0; i < worker; i++ {
		go func(i int) {
			j := 0
			timeoutChannel := time.After(dur)

			for {
				select {
				case <-timeoutChannel:
					atomic.AddInt32(&done, 1)
					wg.Done()
					return
				default:
					f(i, j)
					j++
					atomic.AddInt32(&loop, 1)
				}
			}
		}(i)
	}

	if len(noDebug) == 0 || !noDebug[0] {
		lastLoop := int32(0)
		ticker := timeUtil.NewTicker(time.Second, time.Second, func(now time.Time) {
			_loop := atomic.LoadInt32(&loop)
			_lastLoop := lastLoop
			var tps int32

			tps, lastLoop = _loop-_lastLoop, _loop
			util.Println("    [%s] [produce] loop=%d,lastLoop=%v,  worker=%d/%d, tps=%v", now.Format("15:04:05.000"), _loop, _lastLoop, done, worker, tps)
		})
		defer func() {
			ticker.Trigger()
			ticker.Stop(0)
		}()
	}

	wg.Wait()

	ttl := time.Since(start)
	fmt.Println("----", ttl.String(), loop, loop/int32(ttl.Seconds()))

}

func newTestQueue(name string, args ...interface{}) *Queue {
	pool := channelPool
	opt := Options{NoAutoDeclare: true, NoProduceLog: true, NoConsumeLog: true}
	for _, v := range args {
		switch t := v.(type) {
		case *ChannelPool:
			pool = t
		case Options:
			opt = t
		case *Options:
			opt = *t
		}
	}
	return New(name, pool, opt)
}

func TestChannelClose(t *testing.T) {
	queue := newTestQueue("test")

	err := queue.Do(func(ch *amqp.Channel) error {
		queue.Stop(0)
		go func() {
			queue.channelPool.conn.Close()
			fmt.Println("aaaaaaa")
		}()

		runtimeUtil.Go(3000, func(i int) {
			fmt.Println("--> ", i, ch.Close())
		})
		return nil
	})
	assert.NoError(t, err)
}
