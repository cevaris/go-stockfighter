package main

import (
	"time"
	"fmt"
	"github.com/cevaris/stockfighter"
)

func main() {
	fmt.Println(time.Now().Unix())

	config := stockfighter.InitConfig(".env.yml")

	fmt.Printf("config: %#v\n", config.ApiKey)

	api := stockfighter.InitApi(config, "HAE23155229")

	if value, err := api.HeartBeat(); err == nil {
		fmt.Printf("request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	if value, err := api.VenueHeartBeat("TESTEX"); err == nil {
		fmt.Printf("request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	venuStockQuotes := make(chan *stockfighter.StockQuote, 100)
	defer close(venuStockQuotes)
	go api.VenueTickerTape(venuStockQuotes, "FOKNEX")
	go func(c chan *stockfighter.StockQuote) {
		for stockQuote := range c {
			fmt.Printf("Venue StreamQuote: %#v\n", stockQuote)
			fmt.Printf("%#v\n", stockQuote.QuoteTime)
		}
	}(venuStockQuotes)

	stockQuotes := make(chan *stockfighter.StockQuote, 100)
	defer close(stockQuotes)
	go api.StockTickerTape(stockQuotes, "FOKNEX", "ACH")
	go func(c chan *stockfighter.StockQuote) {
		for stockQuote := range c {
			fmt.Printf("Stock StreamQuote: %#v\n", stockQuote)
			fmt.Printf("%#v\n", stockQuote.QuoteTime)
		}
	}(stockQuotes)

	//venueExecutions := make(chan *stockfighter.Execution, 100)
	//defer close(venueExecutions)
	//go api.VenueExecutions(venueExecutions, "FOKNEX")
	//go func(c chan *stockfighter.Execution) {
	//	for execution := range c {
	//		fmt.Printf("Venue Execution: %#v\n", execution)
	//		fmt.Printf("%#v\n", execution.Order)
	//	}
	//}(venueExecutions)
	//
	//stockExecutions := make(chan *stockfighter.Execution, 100)
	//defer close(stockExecutions)
	//go api.StockExecutions(stockExecutions, "FOKNEX", "ACH")
	//go func(c chan *stockfighter.Execution) {
	//	for execution := range c {
	//		fmt.Printf("Stock Execution: %#v\n", execution)
	//		fmt.Printf("%#v\n", execution.Order)
	//	}
	//}(stockExecutions)

	time.Sleep(10 * time.Second)
}