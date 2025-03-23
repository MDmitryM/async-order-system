package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const (
	OrderTopic    = "orders"
	PaymentTopic  = "payments"
	ShippingTopic = "shipping"
)

// TODO: потенциально отрефакторить, оставив тут только объявление структуры и конструктор,
// а методы раскидать по файлам f.e.: kafkaOrder.go kafka*.go
type OrderMessage struct {
	ID            int32  `json:"id"`
	UserID        int32  `json:"user_id"`
	Total         int32  `json:"total"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	ProductID     int32  `json:"product_id"`
}

type Producer struct {
	syncProducer sarama.SyncProducer
}

func NewSyncProducer(brokers []string) (*Producer, error) {
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.Return.Successes = true
	producerCfg.Producer.RequiredAcks = sarama.WaitForAll
	producerCfg.Producer.Retry.Max = 5
	producerCfg.Producer.Idempotent = false

	producer, err := sarama.NewSyncProducer(brokers, producerCfg)
	if err != nil {
		return nil, err
	}

	return &Producer{syncProducer: producer}, nil
}

func (p *Producer) SendOrder(ctx context.Context, order OrderMessage) error {
	msgBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: OrderTopic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	logrus.Infof("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
