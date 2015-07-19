package main

import (
	"github.com/cryptobase/scraper/model"
	"github.com/cryptobase/scraper/bitfinex"
	"fmt"
	"log"
)

func main() {
	//Load existing trades from file
	existing_trades, err := LoadFromCsv()
	if err != nil {
		log.Fatal(err)
	}

	//Get latest trade
	last_timestamp := 0
	if len(existing_trades) > 0 {
		//last_record := existing_trades[len(existing_trades)-1]
	}

	//Load new trades from api
	new_trades, err := bitfinex.Scrape(last_timestamp)
	if err != nil {
		log.Fatal(err)
	}

	//Append new trades to file
	for _, t := range new_trades {
		fmt.Printf("%d, %d, %s, %4s %12.8f BTC for %7.2f USD/BTC\n", t.Timestamp, t.TradeId, t.Exchange, t.Type, t.Amount, t.Price)
	}
}

func LoadFromCsv() ([]model.Trade, error) {
	return nil, nil
}

func AppendToCsv() ([]model.Trade, error) {
	return nil, nil
}