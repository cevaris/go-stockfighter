package main

import (
	"fmt"
	"github.com/cevaris/stockfighter"
	"time"
	"math/rand"
	"log"
)

var (
	Venue = "YJEX"
	Symbol = "LTH"
	Account = "FAM20791002"

	config = stockfighter.InitConfig(".env.yml", Account)
	api = stockfighter.InitApi(config)
	random = rand.New(rand.NewSource(time.Now().Unix()))
)

var (
	stockQuotes = make(chan *stockfighter.StockQuote, 100)
	executions = make(chan *stockfighter.Execution, 100)

	outstandingShares uint64 = 0
	maxShares = 100000
)

func main() {

	api.IsExchangeHealthy()

	listenToTicker()
	listenToFills()

	for {
		printStatus()
		time.Sleep(3 * time.Second)
	}
	defer close(stockQuotes)
	defer close(executions)
}

func printStatus() {
	status, err := api.StockOrdersAccountStatus(Venue, Symbol)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)
}

func listenToTicker() {
	go api.StockTickerTape(stockQuotes, Venue, Symbol)
	go func(c chan *stockfighter.StockQuote) {
		for q := range c {
			fmt.Println(q.String())
		}
	}(stockQuotes)
}

func listenToFills() {
	go api.StockExecutions(executions, Venue, Symbol)
	go func(c chan *stockfighter.Execution) {
		for e := range c {
			fmt.Print(e.String())
		}
	}(executions)
}