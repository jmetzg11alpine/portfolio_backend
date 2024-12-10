package alpaca_script

import (
	"backend/services/alpaca"
	"backend/services/database"
	"fmt"
	"log"
)

func Run() error {
	if alpaca.IsMarketOpen() {
		log.Printf("market is open")

		// The total amount of money put aside for each ETF
		reserves := database.GetTotalReserves()

		// The amount of money availabe in the alpaca account
		balance := alpaca.CheckAccount()

		// updating the amount of money put aside for each ETF
		etfReserves := database.UpdateReserves(reserves, balance)

		// get the positions of the etfs
		positions, err := alpaca.GetPositions()
		if err != nil {
			log.Printf("Failed to get positions: %v\n", err)
		}
		fmt.Println(positions)

		for _, position := range positions {
			symbol, symbolOk := position["symbol"].(string)
			percentChange, percentChangeOk := position["change_today"].(float64)
			unrealizedPlpc, unrealizedPlpcOk := position["unrealized_plpc"].(float64)
			marketValue, marketValueOk := position["market_value"].(float64)
			reserve, reserveOk := etfReserves[symbol]
			if !symbolOk {
				log.Printf("Invalid symbol type in positioni")
			}
			if !percentChangeOk {
				log.Printf("No percent change for today")
			}
			if !unrealizedPlpcOk {
				log.Printf("Invalid percent change")
			}
			if !marketValueOk {
				log.Printf("Invalid market value")
			}
			if !reserveOk {
				log.Printf("no reserve found for %s", symbol)
			}
			amountSpent := alpaca.InvestInEtf(symbol, percentChange, reserve)
			log.Printf("ETF: %s, percent change: %f, amount spent: %f\n", symbol, percentChange, amountSpent)
			database.UpdateDabase(symbol, amountSpent, percentChange, unrealizedPlpc, marketValue)

		}

	} else {
		log.Println("market is closed")
	}
	return nil
}
