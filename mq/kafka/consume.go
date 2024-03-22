package kafka

import (
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/IBM/sarama"
	"github.com/panjf2000/ants/v2"
	"sync"
)

func (this *Kafka) ensureConsumer() (err error) {
	if this.consumer != nil {
		return
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.consumer != nil {
		return
	}

	this.consumer, err = sarama.NewConsumerFromClient(this.client)
	if err != nil {
		return err
	}
	return nil
}

func (this *Kafka) consumePartition(topic string, partition int32, offset int64, handler func(msg *sarama.ConsumerMessage) error) {
	defer runtimeUtil.Recover()

	if err := this.ensureConsumer(); err != nil {
		this.opt.Logger.Error("Consume(%s, %d, %d) error: %v", topic, partition, offset, err)
		return
	}

	pc, err := this.consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		this.opt.Logger.Error("ConsumePartition(%s, %d, %d) error: %v", topic, partition, offset, err)
		return
	}
	defer func() {
		if err = pc.Close(); err != nil {
			this.opt.Logger.Error("Consume(%s, %d, %d) close error: %v", topic, partition, offset, err)
		}
	}()

	pool, _ := ants.NewPool(this.opt.ConsumeConcurrent)
	defer pool.Release()

	for {
		select {
		case message := <-pc.Messages():
			if err = pool.Submit(func() {
				if err = handler(message); err != nil {
					this.opt.Logger.Error("Consume(%s, %d, %d) message error: %v", topic, partition, offset, err)
				}
			}); err != nil {
				this.opt.Logger.Error("Consume(%s, %d, %d) pool submit error: %v", topic, partition, offset, err)
			}

		case e := <-pc.Errors():
			this.opt.Logger.Error("Consume(%s, %d, %d) error: %v", e.Topic, e.Partition, e.Err)
		}
	}
}

func (this *Kafka) consume(topic string, offset int64, handler func(msg *sarama.ConsumerMessage) error) {
	partitions, err := this.client.Partitions(topic)
	if err != nil {
		this.opt.Logger.Error("Partitions(%s) error: %v", topic, err)
		return
	}

	var wg sync.WaitGroup
	for _, part := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer runtimeUtil.Recover()
			defer wg.Done()
			this.consumePartition(topic, partition, offset, handler)
		}(part)
	}
	wg.Wait()
}

func (this *Kafka) ConsumePartitionNewest(topic string, partition int32, handler func(msg *sarama.ConsumerMessage) error) {
	this.consumePartition(topic, partition, sarama.OffsetNewest, handler)
}

func (this *Kafka) ConsumePartitionOldest(topic string, partition int32, handler func(msg *sarama.ConsumerMessage) error) {
	this.consumePartition(topic, partition, sarama.OffsetOldest, handler)
}

func (this *Kafka) ConsumeNewest(topic string, handler func(msg *sarama.ConsumerMessage) error) {
	this.consume(topic, sarama.OffsetNewest, handler)
}

func (this *Kafka) ConsumeOldest(topic string, handler func(msg *sarama.ConsumerMessage) error) {
	this.consume(topic, sarama.OffsetOldest, handler)
}
