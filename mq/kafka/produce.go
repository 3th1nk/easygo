package kafka

import (
	"github.com/IBM/sarama"
)

func (this *Kafka) ensureProducer() (err error) {
	if this.producer != nil {
		return
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.producer != nil {
		return
	}

	this.producer, err = sarama.NewSyncProducerFromClient(this.client)
	if err != nil {
		return
	}

	return nil
}

func (this *Kafka) Produce(topic, data string) (partition int32, offset int64, err error) {
	if err = this.ensureProducer(); err != nil {
		return
	}
	return this.producer.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(data)})
}

func (this *Kafka) ensureAsyncProducer() (err error) {
	if this.asyncProducer != nil {
		return
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.asyncProducer != nil {
		return
	}

	this.asyncProducer, err = sarama.NewAsyncProducerFromClient(this.client)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case msg := <-this.asyncProducer.Successes():
				this.opt.Logger.Debug("kafka async produce success", msg.Topic, msg.Partition, msg.Offset)
			case e := <-this.asyncProducer.Errors():
				this.opt.Logger.Error("kafka async produce error", e.Msg.Topic, e.Msg.Partition, e.Msg.Offset, e.Err)
			}
		}
	}()

	return nil
}

func (this *Kafka) AsyncProduce(topic, data string) (err error) {
	if err = this.ensureAsyncProducer(); err != nil {
		return
	}
	this.asyncProducer.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(data)}
	return nil
}
