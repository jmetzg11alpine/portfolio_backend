package gov

import (
	"backend/config"
	"context"
	"fmt"
	"log"
	"time"
)

type Entry struct {
	Message string `json:"message"`
}

type ForeignAid struct {
	ID      int
	Country string
	Amount  string
	Year    int
	Lat     float32
	Lng     float32
}

func GetForeignAidData(country, year string) Entry {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	yearQuery := "SELECT * FROM foreign_aid"
	countryQuery := "SELECT * FROM foreign_aid"

	if year != "all" {
		yearQuery += " WHERE year = ?"
	}
	if country != "all" {
		countryQuery += " WHERE country = ?"
	}

	var yearParams []interface{}
	if year != "all" {
		yearParams = append(yearParams, year)
	}
	yearRows, err := config.DB.QueryContext(ctx, yearQuery, yearParams...)
	if err != nil {
		log.Fatalf("Failed to execute year query: %v", err)
	}
	defer yearRows.Close()

	yearData := []ForeignAid{}
	for yearRows.Next() {
		var aid ForeignAid
		if err := yearRows.scan(&aid.ID, &aid.Country, &aid.Amount, &aid.Lat, &aid.Lng); err != nil {
			yearData = append(yearData, aid)
		}
	}

	return Entry{Message: "hello"}
}
