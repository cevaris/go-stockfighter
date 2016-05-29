package main

import (
	"fmt"
	"github.com/cevaris/stockfighter"
	"time"
)

var (
	Account = "MS48925855"
	Venue = "WUTMEX"
	Symbol = "EYNI"
)

var (
	config = stockfighter.InitConfig(".env.yml", Account)
	api = stockfighter.InitApi(config)
	session = stockfighter.InitSession(config, Venue)
)

func main() {
	session.Observe(Symbol)
	for session.LatestQuote == nil || session.LatestQuote.Ask == 0 {
		fmt.Println("waiting for first quote", session.LatestQuote)
		time.Sleep(1 * time.Second)
	}

	soReq := &stockfighter.StockOrderRequest{
		Account: config.Account,
		Venue: Venue,
		Stock: Symbol,
		Price: session.LatestQuote.Ask - 100,
		Qty: 100,
		Direction: stockfighter.DirectionBuy,
		OrderType: stockfighter.OrderMarketK,
	}
	stockOrder, soResErr := api.StockOrder(soReq);
	if soResErr == nil {
		fmt.Printf("stockorder response: %#v\n", stockOrder)
	} else {
		fmt.Println(soResErr)
	}
}