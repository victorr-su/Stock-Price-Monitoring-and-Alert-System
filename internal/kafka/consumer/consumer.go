package consumer

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)
func StartConsumer(){
	// Set up Kafka consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",   // Kafka server address
		"group.id":          "stock-price-consumer", // Consumer group ID
		"auto.offset.reset": "earliest",           // Start from the earliest offset
	}

	// Create a kafka consumer

	consumer, err := kafka.NewConsumer(config)
	if err != nil{
		log.Fatalf("Failed to create consumer: %s\n", err)
	}

	defer consumer.Close()

	//Subscribe to the topic
	err = consumer.Subscribe("stock_prices", nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s\n", err)
	}

	fmt.Println("Consuming messages from Kafka...")

	//consume messages
	for {
		msg, err := consumer.ReadMessage(-1) // -1 waits indefinitely for a message
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v\n", err)
		}
	}

}