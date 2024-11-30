package alpaca

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"
)

type AccountResponse struct {
	Cash string `json:"cash"`
}

func CheckAccount() float64 {
	url := os.Getenv("accountUrl") + "/account"

	req := createRequest("GET", url, nil)
	body, err := sendRequest(req)
	if err != nil {
		log.Printf("Failed to request CheckAccount: %v", err)
	}

	var account AccountResponse
	err = json.Unmarshal(body, &account)
	if err != nil {
		log.Printf("failed ot unmarshal JSON - CheckAccount: %v", err)
	}

	cashFloat, err := strconv.ParseFloat(account.Cash, 64)
	if err != nil {
		log.Printf("Failed to convert cash to float - CheckAccount: %v", err)
	}

	return float64(cashFloat)
}

func InvestInEtf(etf string, percentChange, etfReserves float64) float64 {
	var amountToInvest float64
	if percentChange >= 0 {
		amountToInvest = 0
	} else if percentChange > -0.25 {
		amountToInvest = etfReserves * .05
	} else if percentChange > -0.5 {
		amountToInvest = etfReserves * .1
	} else if percentChange > -.75 {
		amountToInvest = etfReserves * .15
	} else if percentChange > -1 {
		amountToInvest = etfReserves * .2
	} else if percentChange > -1.24 {
		amountToInvest = etfReserves * .25
	} else if percentChange > -1.5 {
		amountToInvest = etfReserves * .3
	} else if percentChange > -1.75 {
		amountToInvest = etfReserves * .35
	} else if percentChange > -2 {
		amountToInvest = etfReserves * .4
	} else if percentChange > -2.25 {
		amountToInvest = etfReserves * .45
	} else if percentChange > -2.5 {
		amountToInvest = etfReserves * .5
	} else if percentChange > -2.75 {
		amountToInvest = etfReserves * .55
	} else if percentChange > -3 {
		amountToInvest = etfReserves * .6
	} else if percentChange > -3.25 {
		amountToInvest = etfReserves * .65
	} else if percentChange > -3.5 {
		amountToInvest = etfReserves * .7
	} else if percentChange > 3.75 {
		amountToInvest = etfReserves * .75
	} else if percentChange > -4 {
		amountToInvest = etfReserves * .8
	} else if percentChange > -4.25 {
		amountToInvest = etfReserves * .85
	} else if percentChange > -4.5 {
		amountToInvest = etfReserves * .9
	} else if percentChange > -4.75 {
		amountToInvest = etfReserves * .95
	} else {
		amountToInvest = etfReserves
	}

	var amountSpent float64

	if amountToInvest > 1.5 {

		amountToInvest = math.Round(amountToInvest*100) / 100

		url := os.Getenv("accountUrl") + "/orders"
		order := map[string]interface{}{
			"symbol":        etf,
			"notional":      amountToInvest,
			"side":          "buy",
			"type":          "market",
			"time_in_force": "day",
		}
		orderData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order data - InvestInEtf: %v", err)
		}
		req := createRequest("POST", url, orderData)
		_, err = sendRequest(req)
		if err != nil {
			log.Printf("Failed to place order for ETF %s - InvestInEtf: %v", etf, err)
		}

		amountSpent = amountToInvest

	} else {
		amountSpent = 0.0
	}

	return amountSpent
}
