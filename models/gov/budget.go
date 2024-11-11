package gov

import (
	"backend/config"
	"context"
	"fmt"
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

var colors = []string{
	"0, 128, 128",  // Teal
	"255, 99, 71",  // Tomato
	"124, 252, 0",  // Lawn Green
	"70, 130, 180", // Steel Blue
	"255, 215, 0",  // Gold
	"0, 191, 255",  // Deep Sky Blue
	"255, 69, 0",   // Orange Red
	"138, 43, 226", // Blue Violet
	"60, 179, 113", // Medium Sea Green
	"218, 165, 32", // Golden Rod
}

type MainEntry struct {
	Label           string  `json:"label"`
	Value           float32 `json:"value"`
	BackgroundColor string  `json:"backgroundColor"`
}

type SimpleEntry struct {
	Label string  `json:"label"`
	Value float32 `json:"value"`
}

func getAgencyAndColor(i int, agency string) (string, string) {
	color := colors[i%9]
	if renamed, exists := renaming[agency]; exists {
		return renamed, color
	}
	return agency, color
}

func makeMainEntry(label, color string, value float32) *MainEntry {
	return &MainEntry{
		Label:           label,
		Value:           value,
		BackgroundColor: fmt.Sprintf("rgba(%s, 1)", color),
	}
}

func GetAgencyData() ([]*MainEntry, []*MainEntry, []*SimpleEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := config.DB.Query(ctx, "SELECT agency, budget FROM agency_budget ORDER BY budget DESC")
	if err != nil {
		log.Printf("Query error: %v\n", err)
		return nil, nil, nil, err
	}
	defer rows.Close()

	var mainData, otherData []*MainEntry
	var remainingData []*SimpleEntry
	var mainOtherValue, otherOtherValue float32
	counter := 0

	for rows.Next() {
		var label string
		var value float32
		if err := rows.Scan(&label, &value); err != nil {
			log.Printf("Row scan error: %v\n", err)
			continue
		}
		shortName, color := getAgencyAndColor(counter, label)
		entry := makeMainEntry(shortName, color, value)

		switch {
		case counter < 9:
			mainData = append(mainData, entry)
		case counter < 18:
			mainOtherValue += value
			otherData = append(otherData, entry)
		default:
			mainOtherValue += value
			otherOtherValue += value
			remainingData = append(remainingData, &SimpleEntry{Label: label, Value: value})
		}

		counter++
	}
	mainData = append(mainData, makeMainEntry(
		"other agencies", colors[9], mainOtherValue,
	))
	otherData = append(otherData, makeMainEntry(
		fmt.Sprintf("%d others", len(remainingData)), colors[9], otherOtherValue,
	))

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v\n", err)
		return nil, nil, nil, err
	}

	return mainData, otherData, remainingData, nil
}
