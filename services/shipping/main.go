package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MDmitryM/async-order-system/services/shipping/kafka"
	"github.com/sirupsen/logrus"
)

var (
	KAFKA_BROKERS = []string{
		"kafka1:29091",
		"kafka2:29092",
		"kafka3:29093",
	}

	CONSUMER_GROUP = "shipping-group"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	fmt.Println("shipping service main")

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shippingChan := make(chan kafka.ShippingMessage, 100)

	producer, err := kafka.NewProducer(KAFKA_BROKERS)
	if err != nil {
		logrus.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(KAFKA_BROKERS, CONSUMER_GROUP, shippingChan)
	if err != nil {
		logrus.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.Run(rootCtx, shippingChan)
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
	close(shippingChan)
	wg.Wait()

	logrus.Info("Shipping service stopped")
}
