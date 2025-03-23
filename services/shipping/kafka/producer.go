package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.Return.Successes = true
	producerCfg.Producer.RequiredAcks = sarama.WaitForAll
	producerCfg.Producer.Retry.Max = 5
	producerCfg.Producer.Idempotent = false

	producer, err := sarama.NewSyncProducer(brokers, producerCfg)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) SendShipping(ctx context.Context, order ShippingMessage) error {
	msgBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: ShippingTopic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	logrus.Infof("Sent shipping message to partition %d at offset %d", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) Run(ctx context.Context, shippingChan <-chan ShippingMessage) {
	for {
		select {
		case <-ctx.Done():
			logrus.Info("Producer stopped by context")
			return
		case shipping, ok := <-shippingChan:
			if !ok {
				logrus.Info("Shipping channel closed")
				return
			}
			if err := p.SendShipping(ctx, shipping); err != nil {
				logrus.Errorf("Failed to send shipping message: %v", err)
			}
		}
	}
}
