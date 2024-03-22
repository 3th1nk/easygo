package kafka

import (
	"github.com/3th1nk/easygo/util/logs"
	"github.com/IBM/sarama"
)

type Options struct {
	ConsumeConcurrent int
	Config            *sarama.Config
	Logger            logs.Logger
}

var DefaultOptions = Options{
	ConsumeConcurrent: 1000,
	Config:            sarama.NewConfig(),
	Logger:            logs.Default,
}
