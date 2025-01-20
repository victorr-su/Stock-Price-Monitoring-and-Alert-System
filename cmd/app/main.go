package main

import (
	"Stock-Price-Monitoring-and-Alert-System/internal/alert"
	"Stock-Price-Monitoring-and-Alert-System/internal/config"
	"Stock-Price-Monitoring-and-Alert-System/internal/kafka/producer"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

// Define a hashmap and mutex to store stock prices and ensure thread safety
var stockPriceMap = make(map[string]float64)
var mutex = &sync.Mutex{}

// EmailData contains the dynamic data to populate the email template
type EmailData struct {
	StockSymbol     string
	PercentageChange float64
	Price           float64
	Time            string
}

func main() {
	fmt.Println("Stock price application started")

	// Load stock configurations
	stockConfigs, err := config.LoadStocksConfig("/root/Stock-Price-Monitoring-and-Alert-System/internal/config/stocks.json")
	if err != nil {
		log.Fatalf("Failed to load stock configurations: %v", err)
	}

	// Initialize Kafka producer
	err = producer.StartProducer()
	if err != nil {
		log.Fatalf("Failed to initialize producer: %v", err)
	}

	// Start a goroutine to clear the hashmap every 24 hours or at a specific time (e.g., market close at 4:00 PM)
	go func() {
		for {
			// Get current time
			now := time.Now()

			// Calculate the next reset time (Market close time - 4:00 PM)
			nextClearTime := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location()) // 4:00 PM local time

			// If the current time is after 4:00 PM, set the next clear time to tomorrow
			if now.After(nextClearTime) {
				nextClearTime = nextClearTime.Add(24 * time.Hour)
			}

			// Sleep until the next scheduled reset time
			time.Sleep(time.Until(nextClearTime))

			// Clear the hashmap at the next reset time
			mutex.Lock()
			stockPriceMap = make(map[string]float64) // Clear the hashmap
			mutex.Unlock()
			log.Println("Cleared stock price map at market close.")
		}
	}()

	// Start the main loop to fetch stock prices and send notifications
	go func() {
		for {
			for _, stock := range stockConfigs {
				// Fetch the stock price
				response, err := producer.FetchStockPrice(stock.Symbol)
				if err != nil {
					log.Printf("Error fetching stock price for %s: %v", stock.Symbol, err)
					continue
				}

				var price float64
				for _, data := range response.TimeSeries {
					closePrice, exists := data["4. close"]
					if !exists {
						log.Printf("Closing price not found for %s", stock.Symbol)
						continue
					}
					price, err = strconv.ParseFloat(closePrice, 64)
					if err != nil {
						log.Printf("Failed to parse price for %s: %v", stock.Symbol, err)
						continue
					}
					break
				}

				// Update the hashmap and check for a price change
				mutex.Lock()
				prevPrice, exists := stockPriceMap[stock.Symbol]
				if exists {
					percentageChange := ((price - prevPrice) / prevPrice) * 100
					// if percentageChange >= 5 || percentageChange <= -5 {
						// Create the email data structure
						emailData := alert.EmailData{
							StockSymbol:     stock.Symbol,
							PercentageChange: percentageChange,
							Price:           price,
							Time:            time.Now().Format("2006-01-02 15:04:05"),
						}

						// Replace with your actual recipient email
						emailRecipient := "su.victor03@gmail.com"
						
						// Send the email using the template
						err = alert.SendEmail(emailRecipient, "Stock Price Alert", emailData)
						if err != nil {
							log.Printf("Failed to send email for %s: %v", stock.Symbol, err)
						} else {
							log.Printf("Notification sent for %s: %.2f%% change.", stock.Symbol, percentageChange)
						}
					// }
				}
				// Update the hashmap with the latest price
				stockPriceMap[stock.Symbol] = price
				mutex.Unlock()

				// Produce the stock price to Kafka
				err = producer.SendStockPrice(stock.Symbol, price)
				if err != nil {
					log.Printf("Failed to send price for %s to Kafka: %v", stock.Symbol, err)
				} else {
					log.Printf("Produced price for %s: %.2f", stock.Symbol, price)
				}
			}
			// Fetch prices every 30 minutes
			time.Sleep(180 * time.Minute)
		}
	}()

	select {} // Prevent the program from exiting
}
