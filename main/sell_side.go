package main

import (
	"fmt"
	"github.com/cevaris/stockfighter"
	"time"
)

var (
	Account = "HB43284293"
	Venue = "QWJEX"
	Symbol = "OMM"
)

var (
	config = stockfighter.InitConfig(".env.yml", Account)
	api = stockfighter.InitApi(config)
	session = stockfighter.InitSession(config, Venue)
)

func main() {
	api.IsExchangeHealthy()
	session.Observe(Symbol)

	go positionWorker()

	for {
		fmt.Println(session)
		time.Sleep(3 * time.Second)
	}
}

func positionWorker() {
	for {
		currQuote, err := api.StockQuote(Venue, Symbol)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if currQuote.Ask > 0 {
			executeOrder(stockfighter.DirectionBuy, currQuote.Ask - 400, 50)
		}
		if currQuote.Bid > 0 {
			executeOrder(stockfighter.DirectionSell, currQuote.Bid + 400, 50)
		}
		time.Sleep(3 * time.Second)
	}
}

func executeOrder(direction string, price int, qty int) *stockfighter.StockOrder {
	soReq := &stockfighter.StockOrderRequest{
		Account: Account,
		Venue: Venue,
		Stock: Symbol,
		Price: price,
		Qty: qty,
		Direction: direction,
		OrderType: stockfighter.OrderLimit,
	}

	fmt.Println("Order", soReq.String())

	so, soResErr := api.StockOrder(soReq);
	if soResErr != nil {
		fmt.Println(soResErr)
	}
	return so
}

