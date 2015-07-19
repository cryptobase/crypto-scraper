package model

type Trade struct {
	Timestamp 	int64
	TradeId		int64
	Price		float64
	Amount		float64
	Exchange	string
	Type		string
}