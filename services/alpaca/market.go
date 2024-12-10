package alpaca

import (
	"encoding/json"
	"log"
	"os"
)

type ClockResponse struct {
	IsOpen bool `json:"is_open"`
}
type CalendarDay struct {
	Date string `json:"date"`
}
type Bar struct {
	Close float64 `json:"c"`
}
type BarResponse struct {
	Bars map[string][]Bar `json:"bars"`
}
type Trade struct {
	Price float64 `json:"p"`
}
type TradeResponse struct {
	Trade Trade `json:"trade"`
}

func IsMarketOpen() bool {
	url := os.Getenv("accountUrl") + "/clock"
	req := createRequest("GET", url, nil)
	body, err := sendRequest(req)
	if err != nil {
		log.Printf("Failed to request IsMarketOpen: %v", err)
	}

	var clockResponse ClockResponse
	err = json.Unmarshal(body, &clockResponse)
	if err != nil {
		log.Printf("Failed to unmarshal JSON - IsMarketOpen: %v", err)
	}

	return clockResponse.IsOpen
}

func GetPositions() ([]map[string]interface{}, error) {
	url := os.Getenv("accountUrl") + "/positions"
	req := createRequest("GET", url, nil)
	body, err := sendRequest(req)
	if err != nil {
		log.Printf("Faile to request GetPositions: %v\n", err)
	}
	var positions []map[string]interface{}
	if err := json.Unmarshal(body, &positions); err != nil {
		log.Printf("Failed to parse JSON: %v\n", err)
		return nil, err
	}

	var result []map[string]interface{}
	for _, position := range positions {
		marketValue, _ := toFloat64(position["market_value"])
		unrealizedPlpc, _ := toFloat64(position["unrealized_plpc"])
		changeToday, _ := toFloat64(position["change_today"])
		extracted := map[string]interface{}{
			"symbol":          position["symbol"],
			"market_value":    marketValue,
			"unrealized_plpc": unrealizedPlpc,
			"change_today":    changeToday,
		}
		result = append(result, extracted)
	}
	return result, nil
}
