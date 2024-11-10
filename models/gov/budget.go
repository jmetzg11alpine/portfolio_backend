package gov

import (
	"backend/config"
	"context"
	"log"
	"time"
)

var renaming = map[string]string{
	"Department of the Treasury":                  "Treasury",
	"Department of Health and Human Services":     "Health and Human",
	"Department of Defense":                       "Defense",
	"Social Security Administration":              "Social Security",
	"Department of Veterans Affairs":              "Veterans Affairs",
	"Department of Agriculture":                   "Agriculture",
	"Office of Personnel Management":              "OPM",
	"Department of Housing and Urban Development": "Housing",
	"Department of Transportation":                "Transportation",
	"Department of Homeland Security":             "Homeland Security",
	"Department of Energy":                        "Energy",
	"Department of Commerce":                      "Commerce",
	"Department of Education":                     "Education",
	"Environmental Protection Agency":             "Environmental",
	"Department of the Interior":                  "Interior",
	"Department of State":                         "State",
	"General Services Administration":             "General Services",
	"Department of Justice":                       "Justice",
	"Department of Labor":                         "Labor",
	"Pension Benefit Guaranty Corporation":        "Pension",
}

func getAgencyName(agency string) string {
	if renamed, exists := renaming[agency]; exists {
		return renamed
	}
	return agency
}

type DataEntry struct {
	Label string  `json:"label"`
	Value float32 `json:"value"`
}

func GetAgencyData() ([], [] DataEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, "SELECT agency, budget FROM agency_budget")
	if err != nil {
		log.Printf("Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var MainData []DataEntry
	var OtherData []DataEntry

	for rows.Next() {
		var label string
		var value float32
		if err := rows.Scan(&label, &value); err != nil {
			log.Printf("Row scan error: %v\n", err)
		}
		entry := MainData{
			Label: getAgencyName(label),
			Value: value,
		}

		data = append(data, entry)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v\n", err)
		return nil, err
	}

	return data, nil
}
