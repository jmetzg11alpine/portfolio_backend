package gov

import (
	"backend/config"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type BarData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func GetForeignAidBarData(country string) *[]BarData {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query, params := buildBarQuery(country)
	rows, err := executeBarQuery(ctx, query, params)
	if err != nil {
		log.Fatalf("Failed to execute bar query: %v", err)
	}

	barData := processBarRows(rows)

	return barData
}

func buildBarQuery(country string) (string, []interface{}) {
	if country == "all" {
		return `
			SELECT year, SUM(amount) as total_amount
			FROM foreign_aid
			GROUP BY year
			ORDER BY year
		`, nil
	} else {
		return `
			SELECT year, amount
			FROM foreign_aid
			WHERE country = $1
			GROUP BY year, amount
			ORDER BY year
		`, []interface{}{country}
	}
}

func executeBarQuery(ctx context.Context, query string, params []interface{}) (pgx.Rows, error) {
	rows, err := config.DB.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func processBarRows(rows pgx.Rows) *[]BarData {
	var data []BarData

	for rows.Next() {
		var aid BarData
		err := rows.Scan(&aid.X, &aid.Y)
		if err != nil {
			log.Printf("Failed to scan row for bar data: %v", err)
			continue
		}
		data = append(data, aid)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("erro iteration over rows for bar data: %v", err)
		return nil
	}
	return &data
}
