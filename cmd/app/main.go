package main

import (
	"Stock-Price-Monitoring-and-Alert-System/internal/config"
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/consumer"
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/producer"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Stock price application started")

	// Load stock configurations from the config file
	stockConfigs, err := config.LoadStocksConfig("/Users/victorsu/Desktop/Stock-Price-Project/Stock-Price-Monitoring-and-Alert-System/internal/config/stocks.json")
	if err != nil {
		log.Fatalf("Failed to load stock configurations: %v", err)
	}

	// Initialize the Kafka producer
	err = producer.StartProducer()
	if err != nil {
		log.Fatalf("Failed to initialize producer: %v", err)
	}

	// Start the Kafka producer in a separate goroutine
	go func() {
		for {
			for _, stock := range stockConfigs {
				// In place of FetchStockPrice, use a static price for testing
				price := 100.0

				// Produce the stock price to Kafka
				err = producer.SendStockPrice(stock.Symbol, price)
				if err != nil {
					log.Printf("Failed to send price for %s to Kafka: %v", stock.Symbol, err)
				} else {
					log.Printf("Produced price for %s: %.2f", stock.Symbol, price)
				}
			}
			// Sleep for a fixed interval before sending the next price (1 minute in this case)
			time.Sleep(1 * time.Second)
		}
	}()

	// Start the Kafka consumer
	consumer.StartConsumer()
}
