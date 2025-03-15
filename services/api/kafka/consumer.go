package kafka

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/sirupsen/logrus"
)

type PaymentMessage struct {
	OrderID int32  `json:"order_id"`
	Status  string `json:"status"`
}

type ShippingMessage struct {
	OrderID int32  `json:"order_id"`
	Status  string `json:"status"`
}

type Consumer struct {
	client  sarama.ConsumerGroup
	groupID string
	repo    repository.Repository
}

func NewConsumer(brokers []string, groupID string, repo repository.Repository) (*Consumer, error) {
	consCfg := sarama.NewConfig()
	consCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	consCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consCfg.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(brokers, groupID, consCfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:  client,
		groupID: groupID,
		repo:    repo,
	}, nil
}

func (c *Consumer) Close() error {
	return c.client.Close()
}

func (c *Consumer) Consume(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	CGHandler := NewConsumerGroupHandler(c.repo)

	for {
		err := c.client.Consume(ctx, []string{OrderTopic, PaymentTopic, ShippingTopic}, CGHandler)
		if err != nil {
			if err == context.Canceled {
				logrus.Info("Consumer stopped by context")
				return nil
			}
			logrus.Errorf("Consumer error: %v", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
