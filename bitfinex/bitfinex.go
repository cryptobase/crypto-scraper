package bitfinex

import (
	"github.com/cryptobase/scraper/model"
	"github.com/cryptobase/scraper/restclient"
	"encoding/json"
	//"fmt"
	"strconv"
)

type BitfinexTrade struct {
	Timestamp 	int64
	Tid			int64
	Price		string
	Amount		string
	Exchange	string
	Type		string
}

//type BitfinexTrades struct {
//	Trades []BitfinexTrade
//}

func Scrape(last_timestamp int32) ([]model.Trade, error) {
	content, err := restclient.GetContent("https://api.bitfinex.com/v1/trades/btcusd")
	if err != nil {
		return nil, err
	}

	//var trades BitfinexTrades
	//err = json.Unmarshal(content, &trades.Trades)
	bitfinex_trades := make([]BitfinexTrade,0)
	err = json.Unmarshal(content, &bitfinex_trades	)
	if err != nil {
		return nil, err
	}

	trades := []model.Trade{}
	for _, trade := range bitfinex_trades {
		var t model.Trade
		t.Timestamp = trade.Timestamp
		t.TradeId = trade.Tid
		t.Exchange = "bitfinex"
		t.Type = trade.Type
		t.Price, err = strconv.ParseFloat(trade.Price, 64)
		t.Amount, err = strconv.ParseFloat(trade.Amount, 64)

		trades = append(trades, t)
	}

	return trades, nil
}