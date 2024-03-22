package rabbitMQ

import (
	"context"
	"github.com/3th1nk/easygo/util/logs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IProducer interface {
	Close()
	Send(exchange string, queue string, data string) error
	SendToExchange(exchange string, data string) error
	SendToQueue(queue string, data string) error
}

func (this *Queue) NewProducer() (IProducer, error) {
	ch, err := this.channelPool.getChannel("producer")
	if err != nil {
		return nil, err
	}

	return &Producer{chanWrap: ch, queue: this}, nil
}

func (this *Queue) NewGroupProducer(size int) (gp IProducer, err error) {
	arr := make([]IProducer, size)

	defer func() {
		if err != nil {
			for _, a := range arr {
				if a != nil {
					a.Close()
				}
			}
		}
	}()

	for i := 0; i < size; i++ {
		arr[i], err = this.NewProducer()
		if err != nil {
			return nil, err
		}
	}

	return &GroupProducer{size: size, arr: arr}, nil
}

type Producer struct {
	queue *Queue
	*chanWrap
}

func (p *Producer) Send(exchange, queue string, data string) error {
	err := p.send(exchange, queue, data)
	if err != nil {
		p.queue.GetLogger().Error("[mq.%s] produce error: %v, exchange=%v, queue=%v, msg=%v", p.queue.name, err, exchange, queue, data)
	} else if !p.queue.opt.NoProduceLog && logs.IsDebugEnable(p.queue.GetLogger()) {
		if exchange != "" {
			p.queue.GetLogger().Debug(`[mq.%s] produce: exchange=%v, msg=%v`, p.queue.name, exchange, data)
		} else {
			p.queue.GetLogger().Debug(`[mq.%s] produce: queue=%v, msg=%v`, p.queue.name, queue, data)
		}
	}
	return nil
}

func (p *Producer) send(exchange, queue string, data string) error {
	// 网络连接可能导致channel为nil的情况，这里进行重连
	if p.channel == nil || p.channel.IsClosed() {
		if err := p.open(); err != nil {
			return err
		}
	}

	err := p.channel.PublishWithContext(context.Background(), exchange, queue,
		false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, Body: []byte(data)})
	if err == nil {
		return nil
	} else if err != amqp.ErrClosed {
		return err
	} else if err := p.open(); err != nil {
		return err
	} else {
		return p.channel.PublishWithContext(context.Background(), exchange, queue,
			false, false, amqp.Publishing{DeliveryMode: amqp.Persistent, Body: []byte(data)})
	}
}

func (p *Producer) SendToExchange(exchange, data string) error {
	return p.Send(exchange, "", data)
}

func (p *Producer) SendToQueue(queue, data string) error {
	return p.Send("", queue, data)
}

func (p *Producer) Close() {
	if p != nil && p.channel != nil {
		p.close()
		p.pool.putChannel(p.chanWrap)
	}
}

type GroupProducer struct {
	size int
	idx  int
	arr  []IProducer
}

func (p *GroupProducer) getIdx() int {
	p.idx++
	idx := p.idx
	if idx >= p.size {
		idx, p.idx = 0, 0
	}
	return idx
}

func (p *GroupProducer) Send(exchange, queue string, data string) error {
	return p.arr[p.getIdx()].Send(exchange, queue, data)
}

func (p *GroupProducer) SendToExchange(exchange string, data string) error {
	return p.arr[p.getIdx()].SendToExchange(exchange, data)
}

func (p *GroupProducer) SendToQueue(queue string, data string) error {
	return p.arr[p.getIdx()].SendToQueue(queue, data)
}

func (p *GroupProducer) Close() {
	for _, a := range p.arr {
		a.Close()
	}
}
