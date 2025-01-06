package main

import (
	"Stock-Price-Monitoring-and-Alert-System/internal/config"
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/consumer"
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/producer"
	"fmt"
	"log"
	"strconv"
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
				// Fetch stock price
				response, err := producer.FetchStockPrice(stock.Symbol)
				if err != nil {
					log.Printf("Error fetching stock price for %s: %v", stock.Symbol, err)
					continue
				}

				var price float64
				// Loop through the response and find the closing price
				for _, data := range response.TimeSeries {
					// Check if close exists in the map
					closePrice, exists := data["4. close"]
					if !exists {
						log.Printf("Closing price not found for %s", stock.Symbol)
						continue
					}
					price, err = strconv.ParseFloat(closePrice, 64)
					if err != nil {
						log.Printf("Failed to convert price for %s: %v", stock.Symbol, err)
						continue
					}
					break // Exit loop once we find the closing price
				}

				// Check if the price was successfully parsed
				if err != nil {
					log.Printf("No valid closing price found for %s", stock.Symbol)
					continue
				}

				// Produce the stock price to Kafka
				err = producer.SendStockPrice(stock.Symbol, price)
				if err != nil {
					log.Printf("Failed to send price for %s to Kafka: %v", stock.Symbol, err)
				} else {
					log.Printf("Produced price for %s: %.2f", stock.Symbol, price)
				}
			}
			time.Sleep(30 * time.Minute)
		}
	}()
	// Start the Kafka consumer
	consumer.StartConsumer()
}
