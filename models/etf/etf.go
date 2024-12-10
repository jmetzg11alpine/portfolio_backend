package etf

import (
	"backend/config"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ETFReserve struct {
	Stock string  `json:"stock"`
	Value float64 `json:"value"`
}

func GetETFReserves() []ETFReserve {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, "SELECT stock, value FROM stock_values")
	if err != nil {
		log.Printf("Query error in GetETFReserves: %v\n", err)
		return nil
	}
	defer rows.Close()

	var reserves []ETFReserve
	for rows.Next() {
		var reserve ETFReserve
		if err := rows.Scan(&reserve.Stock, &reserve.Value); err != nil {
			log.Printf("Row scan error in GetETFReserves: %v", err)
		}
		reserves = append(reserves, reserve)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows interations error in GetETFReserves: %v", err)
	}

	return reserves
}

type ETFTransaction struct {
	Stock   string    `json:"stock"`
	Date    time.Time `json:"date"`
	Percent string    `json:"percent"`
	Value   float64   `json:"value"`
}

func GetETFTransactions() []ETFTransaction {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, "SELECT stock, created_at, percent, value from stock_transactions ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Query error in GetETFTransactions: %v", err)
	}
	defer rows.Close()

	var transactions []ETFTransaction
	for rows.Next() {
		var transaction ETFTransaction
		if err := rows.Scan(&transaction.Stock, &transaction.Date, &transaction.Percent, &transaction.Value); err != nil {
			log.Printf("Row scan error in GetETFTransactions: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows interactions error in GetETFTransactions: %v", err)
	}
	return transactions
}

func createRequest(method, url string, body []byte) *http.Request {
	apiKey := os.Getenv("key")
	secretKey := os.Getenv("secret")

	if apiKey == "" || secretKey == "" {
		log.Printf("Api key or secret key not found in environment")
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID", apiKey)
	req.Header.Add("APCA-API-SECRET-KEY", secretKey)
	req.Header.Add("Content-Type", "application/json")

	return req
}

func GetStocks() {

	url := os.Getenv("accountUrl") + "/account"
	// url := os.Getenv("accountUrl") + "/positions"
	fmt.Printf("url: %s", url)

	req := createRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to request GetStocks: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code GetStocks: %v\n", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
