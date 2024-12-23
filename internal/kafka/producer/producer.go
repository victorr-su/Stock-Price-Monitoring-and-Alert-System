package producer

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)
func StartProducer(){
	// setup kafka to run on port 9092
	config := &kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}

	//setup the producer
	producer, err := kafka.NewProducer(config)
	if err != nil{
		log.Fatalf("Failed to create producer: %s\n", err)
	}
	
	defer producer.Close()

	topic := "stock_prices"

	//example stock price, REPLACE WITH ACTUAL STOCK PRICES
	stockPrice := `{"symbol": "AAPL", "price": 174.75, "timestamp": "2024-12-18T17:30:00"}`

	//produce the message to kafka
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: []byte(stockPrice),
	}, nil)

	// Wait for message to be delivered
	producer.Flush(15 * 1000)
	
	fmt.Println("Message sent to Kafka!")
}