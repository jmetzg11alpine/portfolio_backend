package alpaca_script

import (
	"backend/services/alpaca"
	"backend/services/database"
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
		etfValues := database.UpdateReserves(reserves, balance)

		// get last day market was open
		lastDay := alpaca.GetLastTradingDay()

		// update each eft based on percent change
		for _, etf := range database.EtfList {
			previousPrice := alpaca.GetPreviousClose(etf, lastDay)
			currentPrice := alpaca.GetCurrentPrice(etf)
			percentChange := ((currentPrice - previousPrice) / previousPrice) * 100
			amountSpent := alpaca.InvestInEtf(etf, percentChange, etfValues[etf])
			log.Printf("ETF: %s, previouse: %f, current: %f, percent change: %f, amount spent: %f\n", etf, previousPrice, currentPrice, percentChange, amountSpent)
			database.UpdateDabase(etf, amountSpent, percentChange)
		}

	} else {
		log.Println("market is closed")
	}
	return nil
}
