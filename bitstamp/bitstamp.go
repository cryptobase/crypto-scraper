package bitstamp

import (
	"github.com/cryptobase/scraper/model"
	"github.com/cryptobase/scraper/restclient"
	"encoding/json"
	"strconv"
)


type BitstampTrade struct {
	Date 		string
	Tid			int64
	Price		string
	Amount		string
	Exchange	string
	Type		int
}

func Scrape(last_timestamp uint32) ([]model.Trade, error) {
	url := "https://www.bitstamp.net/api/transactions/" //defaults to all trades from last hour
	content, err := restclient.GetContent(url)
	if err != nil {
		return nil, err
	}

	bitstamp_trades := make([]BitstampTrade,0)
	err = json.Unmarshal(content, &bitstamp_trades	)
	if err != nil {
		return nil, err
	}

	//Bitstamp trade response has latest trade first, reverse
	trades := []model.Trade{}
	last := len(bitstamp_trades)-1
	for i, _ := range bitstamp_trades {
		trade := bitstamp_trades[last-i]

		var t model.Trade
		ts, _ := strconv.ParseUint(trade.Date, 10, 32)
		t.Timestamp = uint32(ts)
		t.TradeId = trade.Tid
		t.Exchange = "bitstamp"
		if trade.Type == 0 {
			t.Type = "buy"
		} else if trade.Type == 1 {
			t.Type = "sell"
		}
		t.Price, _ = strconv.ParseFloat(trade.Price, 64)
		t.Amount, _ = strconv.ParseFloat(trade.Amount, 64)

		if t.Timestamp > last_timestamp {
			trades = append(trades, t)
		}
	}

	return trades, nil
}