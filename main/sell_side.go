package main

import (
	"fmt"
	"github.com/cevaris/stockfighter"
	"time"
	"math/rand"
	"sync/atomic"
)

var (
	venue = "CIMBEX"
	symbol = "XEW"
	account = "SJB21194249"

	config = stockfighter.InitConfig(".env.yml", account)
	api = stockfighter.InitApi(config)
	random = rand.New(rand.NewSource(time.Now().Unix()))
)

var (
	stockQuotes = make(chan *stockfighter.StockQuote, 100)
	outstandingShares uint64 = 0
	maxShares = 100000
)

func main() {

	api.IsExchangeHealthy()

	listenToTicker()

	for {
		printStatus()
		time.Sleep(3 * time.Second)
	}
	close(stockQuotes)
}

func printStatus() {
	currOS := atomic.LoadUint64(&outstandingShares)
	fmt.Println("curOs:", currOS)
}

func listenToTicker() {
	go api.StockTickerTape(stockQuotes, venue, symbol)
	go func(c chan *stockfighter.StockQuote) {
		for q := range c {
			fmt.Println("Quote", q.Ask, q.Bid)
		}
	}(stockQuotes)
}

func movingAverage() int {
	return (newValue*smoothingFactor) + ( workingAverage * ( 1.0 - smoothingFactor) )
}