package rabbitMQ

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/logs"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync/atomic"
	"time"
)

type ConsumeConfig struct {
	// 交换机，队列绑定的交换机，如果队列不存在，重建队列+绑定交换机
	Exchange string
	Queue    string
	// With a prefetch count greater than zero, the server will deliver that many
	// messages to consumers before acknowledgments are received.
	Qos        int
	Concurrent int
}

func (this *Queue) NewSimpleConsumer(cfg *ConsumeConfig) *SimpleConsumer {
	// https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3
	if cfg.Qos <= 0 {
		cfg.Qos = 30
	}

	c := &SimpleConsumer{
		broker: this,
		cfg:    cfg,
	}

	return c
}

type SimpleConsumer struct {
	broker *Queue
	cfg    *ConsumeConfig
	status int32
	stop   func(wait time.Duration) (stopped bool)
	handle func(msg string) (ack bool, err error)
}

func (c *SimpleConsumer) Handle(handle func(data string) (ack bool, err error)) error {
	if c.IsStop() {
		return fmt.Errorf("consumer is stop")
	}

	err := c.broker.DeclareAndBindExchangeQueue(c.cfg.Exchange, c.cfg.Queue)
	if err != nil {
		return err
	}

	stop, err := c.broker.Consume(c.cfg.Queue, c.cfg.Concurrent, func(data string) (ack bool) {
		ack, err := handle(data)
		if err != nil {
			c.broker.GetLogger().Error("消息消费异常：%v, queue=%v", err, c.cfg.Queue, util.ShortStr(data))
		}
		return
	}, ConsumeOptions{
		Exchange: c.cfg.Exchange,
		OnChannel: func(ch *amqp.Channel) {
			if c.cfg.Qos > 0 {
				if err := ch.Qos(c.cfg.Qos, 0, false); err != nil {
					if logger := c.broker.GetLogger(); logs.IsErrorEnable(logger) {
						logger.Error("queue=%v set PrefetchCount=%v error: %v", c.cfg.Qos, c.cfg.Qos, err)
					}
				}
			}
		},
	})

	c.stop = stop
	return err
}

func (c *SimpleConsumer) Stop(wait time.Duration) error {
	if atomic.CompareAndSwapInt32(&c.status, 0, 1) {
		c.stop(wait)
	}
	return nil
}

func (c *SimpleConsumer) IsStop() bool {
	return c.status == 1
}
