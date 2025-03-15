package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/sirupsen/logrus"
)

type ConsumerGroupHandler struct {
	repo repository.Repository
}

func NewConsumerGroupHandler(repo repository.Repository) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		repo: repo,
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

		case PaymentTopic:
			var payment PaymentMessage
			if err := json.Unmarshal(msg.Value, &payment); err != nil {
				logrus.Errorf("Failed to unmarshall payment message: %v", err)
				continue
			}
			logrus.Infof("Received payment: %+v from partition %d, offset %d", payment, msg.Partition, msg.Offset)

			updatedStatusParams := repository.UpdateOrderStatusParams{
				ID:     payment.OrderID,
				Status: payment.Status,
			}
			if _, err := cgh.repo.UpdateOrderStatus(session.Context(), updatedStatusParams); err != nil {
				logrus.Errorf("Failed to update order status: %v", err)
			}

		case ShippingTopic:
			var shipping ShippingMessage
			if err := json.Unmarshal(msg.Value, &shipping); err != nil {
				logrus.Errorf("Failed to unmarshall payment message: %v", err)
				continue
			}
			logrus.Infof("Received shipping: %+v from partition %d, offset %d", shipping, msg.Partition, msg.Offset)

			updatedStatusParams := repository.UpdateOrderStatusParams{
				ID:     shipping.OrderID,
				Status: shipping.Status,
			}
			if _, err := cgh.repo.UpdateOrderStatus(session.Context(), updatedStatusParams); err != nil {
				logrus.Errorf("Failed to update order status: %v", err)
			}
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
