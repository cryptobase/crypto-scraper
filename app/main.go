package main

import (
	"github.com/op/go-logging"
	"github.com/cryptobase/scraper/model"
	"github.com/cryptobase/scraper/bitfinex"
	"github.com/cryptobase/scraper/bitstamp"
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"strconv"
	"time"
	"io/ioutil"
	"regexp"
)

type Day struct {
	year	int
	month 	int
	day 	int
	trades	[]model.Trade
}

var _log = logging.MustGetLogger("example")

func main() {
	handler_wrapper(bitfinex.Scrape, "bitfinex")
	handler_wrapper(bitstamp.Scrape, "bitstamp")
}

func handler_wrapper(f func(int64) ([]model.Trade, error), name string) {
	existing, new, last_timestamp, err := handler(f, name)
	if err != nil {
		log.Printf("[%10s] Update %7s :: msg=[%s]", name, "failed", err)
	} else {
		log.Printf("[%10s] Update %7s :: msg=[#existing=%9d, last ts=%10d, #new=%5d]", name, "success", existing, last_timestamp, new)
	}
}

func handler(f func(int64) ([]model.Trade, error), name string) (int, int, int64, error) {
	path := "/Users/wilelb/crypto-scraper/"

	err := prepare(path)
	if err != nil {
		return 0,0,0,err
	}

	//Load existing trades from file
	latest_file_name, _ := FindLatestCsvFile(path, name)
	file := fmt.Sprintf("%s%s", path, latest_file_name)
	existing_trades, _ := LoadFromCsv(file)

	//Fetch last timestamp
	last_timestamp := int64(0)
	if len(existing_trades) > 0 {
		last_record := existing_trades[len(existing_trades)-1]
		last_timestamp = last_record.Timestamp
	}

	//Load new trades from api
	new_trades, err := f(last_timestamp)
	if err != nil {
		return 0,0,0,err
	}

	//Append new trades to file
	count, err1 := Persist(path, name, last_timestamp, new_trades)
	if err1 != nil {
		return 0,0,0,err1
	}

	return len(existing_trades), count, last_timestamp, nil
}

func prepare(path string) (error) {
	_, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatal(err)
			return err
		} else {
			log.Printf("Created output directory: [%s]", path)
		}
	}

	return nil
}

func FindLatestCsvFile(path string, name string) (string, error) {
	latest := 0
	var latest_file string

	regex, err := regexp.Compile(`(.*)\.(.*)-(.*)-(.*)\.csv`)
	if err != nil {
		return latest_file, err
	}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		groups := regex.FindStringSubmatch(f.Name())
		if groups != nil {
			if groups[1] == name {
				s := fmt.Sprintf("%s%s%s", groups[2], groups[3], groups[4])
				t, _ := strconv.ParseInt(s, 10, 32)
				if int(t) > latest {
					latest_file = f.Name()
					latest = int(t)
				}
			}
		}
	}
	return latest_file, nil
}

func LoadFromCsv(file string) ([]model.Trade, error) {
	log.Printf("Loading from csv file: %s", file)

	trades := []model.Trade{}

	csvfile, err := os.Open(file)
	if err != nil {
		return trades, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return trades, err
	}

	var trade model.Trade
	for _, record := range rawCSVdata {
		trade.Timestamp, _ = strconv.ParseInt(record[0], 10, 64)
		trade.TradeId, _ = strconv.ParseInt(record[1], 10, 64)
		trade.Exchange = record[2]
		trade.Type = record[3]
		trade.Amount, _ = strconv.ParseFloat(record[4], 64)
		trade.Price, _ = strconv.ParseFloat(record[5], 64)
		trades = append(trades, trade)
	}

	return trades, nil
}

func UnixTimestampToDay(unix int64) (Day) {
	unixTs := time.Unix(unix, 0)
	day := Day{0,0,0, []model.Trade{}}
	day.year = unixTs.Year()
	day.month = int(unixTs.Month())
	day.day = unixTs.Day()
	day.trades = make([]model.Trade, 0)
	return day
}

func SameDay(day1 Day, day2 Day) (bool) {
	return day1.year == day2.year && day1.month == day2.month && day1.day == day2.day
}

func PartitionTradesByDay(last_timestamp int64, trades []model.Trade) ([]Day) {
	//slice containing all days
	days := []Day{}

	//initialize the initial day
	var day Day
	if last_timestamp == 0 {
		day = UnixTimestampToDay(trades[0].Timestamp)
	} else {
		day = UnixTimestampToDay(last_timestamp+1)
	}

	log.Printf("Looping over %d trades", len(trades))
	//partition trade data by day
	for _, trade := range trades {
		day = UnixTimestampToDay(trade.Timestamp)
		current_day := UnixTimestampToDay(trade.Timestamp)

		if SameDay(day, current_day) != true {
			days = append(days, day)
			day = current_day
		}
		day.trades = append(day.trades, trade)

		log.Printf("%#v", day)
	}
	days = append(days, day)

	return days
}

func Persist(path string, name string, last_timestamp int64, trades []model.Trade) (int, error) {
	appended := 0

	log.Printf("Incoming: %d trades", len(trades))

	days := PartitionTradesByDay(last_timestamp, trades)
	log.Printf("%3v", days)
	for _, day := range days {
		log.Printf("Appending: %d trades", len(day.trades))
		count, err := AppendToCsv(path, name, day)
		if err != nil {
			log.Printf("Error: %s", err)
		} else {
			appended += count
		}
	}
	return appended, nil
}

func AppendToCsv(path string, name string, day Day) (int, error) {
	appended := 0



	fname := fmt.Sprintf("%s.%d-%d-%d.csv", name, day.year, day.month, day.day)
	//Check path ends with '/' or look for golang way to resolve paths
	file := fmt.Sprintf("%s%s", path, fname)

	csvfile, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return appended, err
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)
	for _, trade := range day.trades {
		record := []string{
			fmt.Sprintf("%d", trade.Timestamp),
			fmt.Sprintf("%d", trade.TradeId),
			trade.Exchange,
			trade.Type,
			fmt.Sprintf("%.8f", trade.Amount),
			fmt.Sprintf("%.2f", trade.Price)}
		err := writer.Write(record)
		appended++
		if err != nil {
			return appended, err
		}
	}
	writer.Flush()

	return appended, nil
}