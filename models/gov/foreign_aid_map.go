package gov

import (
	"backend/config"
	"context"
	"log"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
)

type Entry struct {
	Message string `json:"message"`
}

type MapData struct {
	Country string  `json:"country"`
	Amount  float64 `json:"amount"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}

func GetForeignAidMapData(year string) (*[]MapData, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query, params := buildMapQuery(year)
	rows, err := executeMapQuery(ctx, query, params)
	if err != nil {
		log.Fatalf("Failed to execute map query: %v", err)
		return nil, nil
	}
	defer rows.Close()

	mapData, countries := processMapRows(rows)

	return mapData, countries
}

func buildMapQuery(year string) (string, []interface{}) {
	if year == "all" {
		return `
			SELECT country, SUM(amount) as total_amount, lat, lng
			FROM foreign_aid
			GROUP BY country, lat, lng
		`, nil
	} else {
		return `
			SELECT country, amount, lat, lng
			FROM foreign_aid
			WHERE year = $1
		`, []interface{}{year}
	}
}

func executeMapQuery(ctx context.Context, query string, params []interface{}) (pgx.Rows, error) {
	rows, err := config.DB.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func processMapRows(rows pgx.Rows) (*[]MapData, []string) {
	var data []MapData
	countrySet := make(map[string]bool)

	for rows.Next() {
		var aid MapData
		err := rows.Scan(&aid.Country, &aid.Amount, &aid.Lat, &aid.Lng)
		if err != nil {
			log.Printf("Falied to scan row: %v", err)
			continue
		}
		data = append(data, aid)
		countrySet[aid.Country] = true
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
		return nil, nil
	}

	var countries []string
	for country := range countrySet {
		countries = append(countries, country)
	}
	sort.Strings(countries)
	countries = append([]string{"all"}, countries...)

	return &data, countries
}
