package etf

import (
	"backend/config"
	"context"
	"fmt"
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

func GetStocks() {
	fmt.Println("i was called")
}
