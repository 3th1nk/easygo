package rabbitMQ

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChannelPool_Publish(t *testing.T) {
	const queueName = "test_0"

	defer func() {
		channelPool.Do("", func(ch *amqp.Channel) {
			_, err := ch.QueueDelete(queueName, false, false, false)
			assert.NoError(t, err)
		})
	}()

	start, total := time.Now(), 100000
	channelPool.Do("", func(ch *amqp.Channel) {
		for i := 0; i < total; i++ {
			err := ch.Publish("", queueName, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         []byte(fmt.Sprintf("test-channel: %d", i)),
			})
			assert.NoError(t, err)
		}
	})
	d := time.Since(start)
	util.Println("total: %d, took: %v, qps=%.2f", total, d.Round(time.Millisecond), float64(total)/d.Seconds())
}
