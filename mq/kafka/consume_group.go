package kafka

import (
	"context"
	"github.com/3th1nk/easygo/util/logs"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/IBM/sarama"
	"github.com/panjf2000/ants/v2"
)

var _ sarama.ConsumerGroupHandler = (*consumeGroupHandler)(nil)

type consumeGroupHandler struct {
	concurrent int
	handler    ConsumeHandler
	logger     logs.Logger
}

func newConsumeGroupHandler(concurrent int, handler ConsumeHandler, logger logs.Logger) *consumeGroupHandler {
	return &consumeGroupHandler{concurrent: concurrent, handler: handler, logger: logger}
}

func (h *consumeGroupHandler) Setup(s sarama.ConsumerGroupSession) error   { return nil }
func (h *consumeGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }
func (h *consumeGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	pool, _ := ants.NewPool(h.concurrent)
	defer pool.Release()

	for msg := range claim.Messages() {
		if err := pool.Submit(func() {
			if err := h.handler(msg); err != nil {
				h.logger.Error("ConsumeGroup message error: %v", err)
			}
			sess.MarkMessage(msg, "")
		}); err != nil {
			h.logger.Error("ConsumeGroup submit error: %v", err)
		}
	}
	return nil
}

func (this *Kafka) ConsumeGroup(ctx context.Context, group string, topic []string, handler ConsumeHandler) {
	defer runtimeUtil.Recover()

	if err := this.ensureConsumer(); err != nil {
		this.opt.Logger.Error("ConsumeGroup(%s) error: %v", group, err)
		return
	}

	cg, err := sarama.NewConsumerGroupFromClient(group, this.client)
	if err != nil {
		this.opt.Logger.Error("ConsumeGroup(%s) error: %v", group, err)
		return
	}
	defer func() {
		if err = cg.Close(); err != nil {
			this.opt.Logger.Error("ConsumeGroup(%s) close error: %v", group, err)
		}
	}()

	go func() {
		defer runtimeUtil.Recover()
		for err = range cg.Errors() {
			this.opt.Logger.Error("ConsumeGroup(%s) error: %v", group, err)
		}
	}()

	for {
		groupHandler := newConsumeGroupHandler(this.opt.ConsumeConcurrent, handler, this.opt.Logger)
		if err = cg.Consume(ctx, topic, groupHandler); err != nil {
			this.opt.Logger.Error("ConsumeGroup(%s, %s) consume error: %v", topic, group, err)
			break
		}
	}
}
