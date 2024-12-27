package producer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer

// Initialize the Kafka producer
func StartProducer() error {
	// Setup Kafka to run on port 9092
	config := &kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}

	// Setup the producer
	var err error
	producer, err = kafka.NewProducer(config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}

	// Return nil without closing the producer here
	return nil
}

// SendStockPrice sends a stock price message to Kafka
func SendStockPrice(symbol string, price float64) error {
	// Check if producer is initialized
	if producer == nil {
		return fmt.Errorf("producer is not initialized")
	}

	// Create the stock price message as JSON
	message := fmt.Sprintf(`{"symbol": "%s", "price": %.2f, "timestamp": "%s"}`,
		symbol, price, time.Now().Format(time.RFC3339))

	// Produce the message to Kafka topic
	topic := "stock-prices"
	deliveryChan := make(chan kafka.Event)
	
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	// Wait for message delivery and handle errors
	e := <-deliveryChan
	msg, ok := e.(*kafka.Message)
	if ok && msg.TopicPartition.Error != nil {
		return fmt.Errorf("failed to deliver message: %v", msg.TopicPartition.Error)
	}

	// Close the delivery channel
	close(deliveryChan)

	// Optionally flush producer if needed for final message delivery
	producer.Flush(15 * 1000)

	fmt.Println("Message sent to Kafka!")
	return nil
}

// CloseProducer gracefully shuts down the producer
func CloseProducer() {
	if producer != nil {
		producer.Close()
	}
}
