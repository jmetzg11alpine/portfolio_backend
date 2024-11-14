package gov

import (
	"backend/config"
	"context"
	"log"
	"sort"
	"time"
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

	mapData, countries := getCountryData(year, ctx)

	return mapData, countries
}

func getCountryData(year string, ctx context.Context) (*[]MapData, []string) {
	if year == "all" {
		return getCountryDataGroupedByCountry(ctx)
	} else {
		return getCountryDataByYear(year, ctx)
	}
}

func getCountryDataGroupedByCountry(ctx context.Context) (*[]MapData, []string) {
	query := `
		SELECT country, SUM(amount) as total_amount, lat, lng
		FROM foreign_aid
		GROUP BY country, lat, lng
	`

	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		log.Fatalf("Failed to exectue grouped by country query %v", err)
		return nil, nil
	}
	defer rows.Close()

	var data []MapData
	countrySet := make(map[string]bool)

	for rows.Next() {
		var aid MapData
		err := rows.Scan(&aid.Country, &aid.Amount, &aid.Lat, &aid.Lng)
		if err != nil {
			log.Printf("Failed to scan row in grouped data: %v", err)
			continue
		}
		data = append(data, aid)
		countrySet[aid.Country] = true
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iteration over rows in grouped data: %v", err)
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

func getCountryDataByYear(year string, ctx context.Context) (*[]MapData, []string) {
	query := `
		SELECT country, amount, lat, lng
		FROM foreign_aid
		WHERE year = $1
	`
	var params []interface{}
	params = append(params, year)

	rows, err := config.DB.Query(ctx, query, params...)
	if err != nil {
		log.Fatalf("Failed to execute country query by year %v:", err)
		return nil, nil
	}
	defer rows.Close()

	var data []MapData
	countrySet := make(map[string]bool)

	for rows.Next() {
		var aid MapData
		err := rows.Scan(&aid.Country, &aid.Amount, &aid.Lat, &aid.Lng)
		if err != nil {
			log.Printf("Failed to scan row in country data: %v", err)
			continue
		}
		data = append(data, aid)
		countrySet[aid.Country] = true
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iteration rows: %v", err)
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
