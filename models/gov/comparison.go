package gov

import (
	"backend/config"
	"context"
	"log"
	"sort"
	"time"
)

type AgencySpending struct {
	Year   int     `json:"year"`
	Amount float64 `json:"amount"`
}

func GetComparisonData() (map[string][]AgencySpending, []int, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, "SELECT year, name, amount from function_spending")
	if err != nil {
		log.Printf("Query error :%v\n", err)
		return nil, nil, nil
	}
	defer rows.Close()

	agencies := make(map[string]struct{})
	years := make(map[int]struct{})
	data := make(map[string][]AgencySpending)

	for rows.Next() {
		var name string
		var year int
		var amount float64
		if err := rows.Scan(&year, &name, &amount); err != nil {
			log.Printf("Row scan error: %v\n", err)
			continue
		}
		agencies[name] = struct{}{}
		years[year] = struct{}{}
		data[name] = append(data[name], AgencySpending{
			Year:   year,
			Amount: amount / 1000000000,
		})
	}

	for key := range data {
		sort.Slice(data[key], func(i, j int) bool {
			return data[key][i].Year < data[key][j].Year
		})
	}

	agencyList := make([]string, 0, len(agencies))
	for agency := range agencies {
		agencyList = append(agencyList, agency)
	}
	sort.Strings(agencyList)

	yearList := make([]int, 0, len(years))
	for year := range years {
		yearList = append(yearList, year)
	}
	sort.Ints(yearList)

	return data, yearList, agencyList
}
