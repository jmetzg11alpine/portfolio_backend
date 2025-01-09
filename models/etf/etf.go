package etf

import (
	"backend/config"
	"context"
	"log"
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
	Id             int       `json:"id"`
	Stock          string    `json:"stock"`
	Date           time.Time `json:"date"`
	Percent        string    `json:"percent"`
	Value          float64   `json:"value"`
	MarketValue    float64   `json:"marketValue"`
	UnrealizedPlpc float64   `json:"unrealizedPlpc"`
}

func GetETFTransactions() []ETFTransaction {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, "SELECT * from stock_transactions WHERE created_at > '2024-12-016' ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Query error in GetETFTransactions: %v", err)
	}
	defer rows.Close()

	var transactions []ETFTransaction
	for rows.Next() {
		var transaction ETFTransaction
		if err := rows.Scan(&transaction.Id, &transaction.Stock, &transaction.Date, &transaction.Percent, &transaction.Value, &transaction.MarketValue, &transaction.UnrealizedPlpc); err != nil {
			log.Printf("Row scan error in GetETFTransactions: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows interactions error in GetETFTransactions: %v", err)
	}
	return transactions
}
