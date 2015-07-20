package bitfinex

import (
	"github.com/cryptobase/scraper/model"
	"github.com/cryptobase/scraper/restclient"
	"encoding/json"
	"strconv"
	"fmt"
)

type BitfinexTrade struct {
	Timestamp 	int64
	Tid			int64
	Price		string
	Amount		string
	Exchange	string
	Type		string
}

func Scrape(last_timestamp uint32) ([]model.Trade, error) {
	url := fmt.Sprintf("https://api.bitfinex.com/v1/trades/btcusd?timestamp=%d", last_timestamp+1)
	content, err := restclient.GetContent(url)
	if err != nil {
		return nil, err
	}

	bitfinex_trades := make([]BitfinexTrade,0)
	err = json.Unmarshal(content, &bitfinex_trades	)
	if err != nil {
		return nil, err
	}

	//Bitfinex trade response has latest trade first, reverse
	trades := []model.Trade{}
	last := len(bitfinex_trades)-1
	for i, _ := range bitfinex_trades {
		trade := bitfinex_trades[last-i]

		var t model.Trade
		t.Timestamp = uint32(trade.Timestamp)
		t.TradeId = trade.Tid
		t.Exchange = "bitfinex"
		t.Type = trade.Type
		t.Price, _ = strconv.ParseFloat(trade.Price, 64)
		t.Amount, _ = strconv.ParseFloat(trade.Amount, 64)

		trades = append(trades, t)
	}

	return trades, nil
}