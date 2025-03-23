package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerGroupHandler struct {
	output chan<- PaymentMessage
}

func NewConsumerGroupHandler(output chan<- PaymentMessage) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		output: output,
	}
}

func (cgh *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		switch msg.Topic {
		case OrderTopic:
			var order OrderMessage
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				logrus.Errorf("Failed to unmarshall order message: %v", err)
				continue
			}
			logrus.Infof("Received order: %+v from partition %d, offset %d", order, msg.Partition, msg.Offset)
			paymentMessage := PaymentMessage{
				OrderID: order.ID,
				Status:  "Paid",
			}

			select {
			case cgh.output <- paymentMessage:
			case <-session.Context().Done():
				logrus.Info("Consumer stopped by context")
				return nil
			}
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
