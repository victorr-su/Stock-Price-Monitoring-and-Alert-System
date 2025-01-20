package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Stock struct {
	Symbol     string  `json:"symbol"`
	TargetPrice float64 `json:"target_price"`
}

type StocksFile struct {
	Stocks []Stock `json:"stocks"`
}

type AlphaVantageResponse struct {
	TimeSeries map[string]map[string]string `json:"Time Series (15min)"`
}

func FetchStockPrice(symbol string) (*AlphaVantageResponse, error){

	// Load the .env file
	err := godotenv.Load("/root/Stock-Price-Monitoring-and-Alert-System")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY");

	if apiKey == ""{
		return nil, fmt.Errorf("API key not set")
	}

	// build the request url
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=%s&interval=15min&apikey=%s", symbol, apiKey)

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error fetching stock price for %s: %v", symbol, err)
	}

	defer resp.Body.Close()

	//parse the response
	var response AlphaVantageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response for %s: %v", symbol, err)
	}

	return &response, nil
}