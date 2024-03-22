package kafka

import (
	"github.com/3th1nk/easygo/util/logs"
	"github.com/IBM/sarama"
	"github.com/modern-go/reflect2"
	"sync"
)

type Kafka struct {
	client        sarama.Client
	consumer      sarama.Consumer
	producer      sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	mu            sync.Mutex
	opt           Options
}

func New(addr []string, opt ...Options) (*Kafka, error) {
	var obj Kafka
	if len(opt) == 0 {
		obj.opt = DefaultOptions
	} else {
		obj.opt = opt[0]
		if obj.opt.Logger == nil {
			obj.opt.Logger = logs.Default
		}
		if obj.opt.Config == nil {
			obj.opt.Config = sarama.NewConfig()
		}
		if obj.opt.ConsumeConcurrent <= 0 {
			obj.opt.ConsumeConcurrent = DefaultOptions.ConsumeConcurrent
		}
	}
	obj.opt.Config.Producer.Return.Successes = true
	obj.opt.Config.Producer.Return.Errors = true

	var err error
	obj.client, err = sarama.NewClient(addr, obj.opt.Config)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (this *Kafka) Close() {
	if !reflect2.IsNil(this.consumer) {
		_ = this.consumer.Close()
	}
	if !reflect2.IsNil(this.producer) {
		_ = this.producer.Close()
	}
	if !reflect2.IsNil(this.asyncProducer) {
		_ = this.asyncProducer.Close()
	}
	if !reflect2.IsNil(this.client) {
		_ = this.client.Close()
	}
}
