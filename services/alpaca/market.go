package alpaca

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
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

func GetLastTradingDay() (lastDay string) {
	url := os.Getenv("accountUrl") + "/calendar"
	req := createRequest("GET", url, nil)
	body, err := sendRequest(req)
	if err != nil {
		log.Printf("Failed to request getLastTradingDay: %v", err)
	}

	var calendar []CalendarDay
	err = json.Unmarshal(body, &calendar)
	if err != nil {
		log.Printf("Failed to unmarshal JSON - GetLastTradingDay: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	for i := len(calendar) - 1; i >= 0; i-- {
		if calendar[i].Date < today {
			lastDay = calendar[i].Date
			return
		}
	}
	log.Printf("No Trading day found before today - GetLatTradingDay")
	return ""
}

func GetPreviousClose(etf, date string) float64 {
	baseUrl := os.Getenv("marketUrl") + "/stocks/bars"
	url := fmt.Sprintf("%s?symbols=%s&timeframe=1Day&start=%s&end=%s", baseUrl, etf, date, date)

	req := createRequest("GET", url, nil)
	body, err := sendRequest(req)
	if err != nil {
		log.Printf("Falied to request GetPreviousClose: %v", err)
	}

	var barResponse BarResponse
	err = json.Unmarshal(body, &barResponse)
	if err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
	}

	if bars, ok := barResponse.Bars[etf]; ok && len(bars) > 0 {
		return bars[0].Close
	}

	log.Printf("No closing data found for %s on %s", etf, date)
	return 0.0
}

func GetCurrentPrice(etf string) float64 {
	baseUrl := os.Getenv("marketUrl") + "/stocks"
	url := fmt.Sprintf("%s/%s/trades/latest", baseUrl, etf)

	req := createRequest("GET", url, nil)
	body, err := sendRequest(req)
	if err != nil {
		log.Printf("Failed to request GetCurrentPrice: %v", err)
	}

	var tradeResponse TradeResponse
	err = json.Unmarshal(body, &tradeResponse)
	if err != nil {
		log.Printf("Failed to unmarshal GetCurrentPrice: %v", err)
	}

	return tradeResponse.Trade.Price
}
