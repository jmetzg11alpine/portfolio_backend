package database

import (
	"backend/config"
	"context"
	"log"
	"math"
)

func GetTotalReserves() float64 {
	ctx := context.Background()
	query := "SELECT value FROM stock_values"

	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	var total float64 = 0.0
	for rows.Next() {
		var value float64

		err := rows.Scan(&value)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
		}
		total += value
	}

	if rows.Err() != nil {
		log.Printf("Error occurred during row iteration: %v", rows.Err())
	}

	return total
}

func UpdateReserves(reserves, balance float64) map[string]float64 {
	realBalance := balance - reserves
	availableFunds := Min(55, .05*realBalance)

	amountEachETF := availableFunds / 11
	amountEachETF = math.Round(amountEachETF*100) / 100

	ctx := context.Background()

	insertQuery := `
		INSERT INTO stock_values (stock, value)
		VALUES($1, $2)
		ON CONFLICT (stock)
		DO UPDATE SET value = stock_values.value + $2
	`
	selectQuery := `
		SELECT value FROM stock_values WHERE stock = $1
	`
	etfValues := make(map[string]float64)

	for _, etf := range EtfList {
		_, err := config.DB.Exec(ctx, insertQuery, etf, amountEachETF)
		if err != nil {
			log.Printf("Failed to update value for stock %s: %v", etf, err)
			continue
		}

		var updatedValue float64
		err = config.DB.QueryRow(ctx, selectQuery, etf).Scan(&updatedValue)
		if err != nil {
			log.Printf("Failed to get updated valuefor stock %s: %v", etf, err)
			continue
		}
		etfValues[etf] = updatedValue
	}

	return etfValues
}

func UpdateDabase(etf string, amountSpent, percentChange, unrealizedPlpc, marketValue float64) {
	ValueQuery := `
		UPDATE stock_values
		SET value = value - $1
		WHERE stock = $2
	`
	ctx := context.Background()
	_, err := config.DB.Exec(ctx, ValueQuery, amountSpent, etf)
	if err != nil {
		log.Printf("Falied to update value for stock %s: %v", etf, err)
	}

	recordQuery := `
		INSERT INTO stock_transactions (stock, percent, value, market_value, unrealized_plpc)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = config.DB.Exec(ctx, recordQuery, etf, percentChange, amountSpent, marketValue, unrealizedPlpc)
	if err != nil {
		log.Printf("Failed to update stock_transactions  in UpdateDatabase eft: %s, %v", etf, err)
	}
}

func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

var EtfList = []string{"IXJ", "RXI", "KXI", "IXG", "EXI", "IXN", "IXC", "MXI", "JXI", "REET", "IXP"}
