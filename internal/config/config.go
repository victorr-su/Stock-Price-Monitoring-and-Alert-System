package config

import (
	"encoding/json"
	"fmt"
	"os"
)
type StockConfig struct {
	Symbol      string  `json:"symbol"`
}

func LoadStocksConfig(filePath string) ([]StockConfig, error){
	file, err := os.Open(filePath)

	if err != nil{
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var stocks []StockConfig
	err = json.NewDecoder(file).Decode(&stocks)

	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	return stocks, nil
}