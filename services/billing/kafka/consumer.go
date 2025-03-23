package kafka

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const (
	OrderTopic    = "orders"
	PaymentTopic  = "payments"
	ShippingTopic = "shipping"
)

type OrderMessage struct {
	ID            int32  `json:"id"`
	UserID        int32  `json:"user_id"`
	Total         int32  `json:"total"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	ProductID     int32  `json:"product_id"`
}

type PaymentMessage struct {
	OrderID int32  `json:"order_id"`
	Status  string `json:"status"`
}

type Consumer struct {
	client  sarama.ConsumerGroup
	groupID string
	output  chan<- PaymentMessage
}

func NewConsumer(brokers []string, groupID string, output chan<- PaymentMessage) (*Consumer, error) {
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
		output:  output,
	}, nil
}

func (c *Consumer) Close() error {
	return c.client.Close()
}

func (c *Consumer) Consume(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	cgHandler := NewConsumerGroupHandler(c.output)

	for {
		err := c.client.Consume(ctx, []string{OrderTopic}, cgHandler)
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
