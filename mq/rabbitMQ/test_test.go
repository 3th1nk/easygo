package rabbitMQ

import (
	"github.com/3th1nk/easygo/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ensureChannelPool()

	m.Run()
}

var channelPool *ChannelPool

func ensureChannelPool() (err error) {
	if channelPool == nil {
		channelPool, err = NewChannelPool("amqp://gsuser:Gs@mq?>>@192.168.1.163:20002")
		return err
	}
	return nil
}

func declareQueue(t *testing.T, queue ...string) {
	channelPool.Do("", func(ch *amqp.Channel) {
		for _, v := range queue {
			_, err := ch.QueueDeclare(v, true, false, false, false, map[string]interface{}{
				"x-expires":    DefaultQueueDeclareOptions.ExpireSec * 1000,
				"x-max-length": DefaultQueueDeclareOptions.Capacity,
				"x-overflow":   DefaultQueueDeclareOptions.Overflow,
			})
			assert.NoError(t, err)
		}
	})
}

func removeQueue(t *testing.T, queue ...string) {
	util.Println("    [%s] begin remove queue", time.Now().Format("15:04:05.000"))
	channelPool.Do("", func(ch *amqp.Channel) {
		for _, v := range queue {
			ch.QueueDelete(v, false, false, false)
		}
	})
	util.Println("    [%s] end remove queue", time.Now().Format("15:04:05.000"))
}
