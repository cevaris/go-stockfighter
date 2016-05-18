package main

import (
	"fmt"
	"github.com/cevaris/stockfighter"
	"log"
)

func main() {
	config := stockfighter.InitConfig(".env.yml")
	fmt.Printf("config: %#v\n", config.ApiKey)

	api := stockfighter.InitApi(config, "HAE23155229")

	if value, err := api.HeartBeat(); err == nil {
		fmt.Printf("request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	stockQuote, err := api.StockQuote("EMCTEX", "TBM")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("current quote: %#v\n", stockQuote)
	}

	soReq := &stockfighter.StockOrderRequest{
		Account: config.Account,
		Venue: stockQuote.Venue,
		Stock: stockQuote.Symbol,
		Price: stockQuote.Bid + 10,
		Qty: 100,
		Direction: "buy",
		OrderType: stockfighter.OrderMarket,
	}
	stockOrder, soResErr := api.StockOrder(soReq);
	if soResErr == nil {
		fmt.Printf("stockorder response: %#v\n", stockOrder)
	} else {
		fmt.Println(soResErr)
	}
}