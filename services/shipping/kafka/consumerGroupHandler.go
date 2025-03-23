package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerGroupHandler struct {
	output chan<- ShippingMessage
}

func NewConsumerGroupHandler(output chan<- ShippingMessage) *ConsumerGroupHandler {
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
		case PaymentTopic:
			var payment PaymentMessage
			if err := json.Unmarshal(msg.Value, &payment); err != nil {
				logrus.Errorf("Failed to unmarshall payment message: %v", err)
				continue
			}
			logrus.Infof("Received payment: %+v from partition %d, offset %d", payment, msg.Partition, msg.Offset)
			shippingMessage := ShippingMessage{
				OrderID: payment.OrderID,
				Status:  "Shipped",
			}

			select {
			case cgh.output <- shippingMessage:
			case <-session.Context().Done():
				logrus.Info("Consumer stopped by context")
				return nil
			}
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
