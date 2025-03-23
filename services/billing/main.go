package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MDmitryM/async-order-system/services/billing/kafka"
	"github.com/sirupsen/logrus"
)

var (
	KAFKA_BROKERS = []string{
		"kafka1:29091",
		"kafka2:29092",
		"kafka3:29093",
	}

	CONSUMER_GROUP = "billing-group"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	fmt.Println("billing service")

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paymentChan := make(chan kafka.PaymentMessage, 100)

	producer, err := kafka.NewProducer(KAFKA_BROKERS)
	if err != nil {
		logrus.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(KAFKA_BROKERS, CONSUMER_GROUP, paymentChan)
	if err != nil {
		logrus.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.Run(rootCtx, paymentChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Consume(rootCtx, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	close(paymentChan)
	wg.Wait()

	logrus.Info("Billing service stopped")
}
