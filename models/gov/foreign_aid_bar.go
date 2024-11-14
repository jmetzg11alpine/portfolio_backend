package gov

import (
	"backend/config"
	"context"
	"log"
	"time"
)

type BarData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func GetForeignAidBarData(country string) *[]BarData {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	barData := getBarData(country, ctx)

	return barData
}

func getBarData(country string, ctx context.Context) *[]BarData {
	if country == "all" {
		return getAllData(ctx)
	} else {
		return getSpecificData(country, ctx)
	}
}

func getAllData(ctx context.Context) *[]BarData {
	query := `
		SELECT year, SUM(amount) as total_amount
		FROM foreign_aid
		GROUP BY year
		ORDER BY year
	`
	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to exectue bar data all countries query: %v", err)
		return nil
	}
	defer rows.Close()

	var data []BarData

	for rows.Next() {
		var aid BarData
		err := rows.Scan(&aid.X, &aid.Y)
		if err != nil {
			log.Printf("Failed to scan row in bar data all countries: %v", err)
			continue
		}
		data = append(data, aid)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("error iteration over rows in bar data all countries: %v", err)
		return nil
	}
	return &data
}

func getSpecificData(country string, ctx context.Context) *[]BarData {
	query := `
		SELECT year, amount
		FROM foreign_aid
		WHERE country = $1
		GROUP BY year, amount
		ORDER BY year
	`
	var params []interface{}
	params = append(params, country)
	rows, err := config.DB.Query(ctx, query, params...)
	if err != nil {
		log.Fatalf("Failed to exectue bar data all countries query: %v", err)
		return nil
	}
	defer rows.Close()

	var data []BarData

	for rows.Next() {
		var aid BarData
		err := rows.Scan(&aid.X, &aid.Y)
		if err != nil {
			log.Printf("Failed to scan row in bar data all countries: %v", err)
			continue
		}
		data = append(data, aid)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("error iteration over rows in bar data all countries: %v", err)
		return nil
	}
	return &data
}
