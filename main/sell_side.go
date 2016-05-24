package main

import (
	"fmt"
	"github.com/cevaris/stockfighter"
	"github.com/cevaris/stockfighter/algo"
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
	smaTri = algo.InitSmaTriple(5, 10, 13)
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
		quote := session.LatestQuote
		smaTri.Push(quote.Last)

		if smaTri.Signal() == stockfighter.SignalBuy {
			if quote.Ask > 0 {
				executeOrder(stockfighter.DirectionBuy, quote.Ask - 100, 50)
			}
		}

		if smaTri.Signal() == stockfighter.SignalBuy {
			if quote.Bid > 0 {
				executeOrder(stockfighter.DirectionSell, quote.Bid + 100, 50)
			}
		}

		time.Sleep(1 * time.Second)
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

