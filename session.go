package stockfighter

import (
	"sync"
	"fmt"
	"encoding/json"
)

type Session struct {
	Cash        int
	NAV         int
	Position    int
	Venue       string
	api         *Api
	config      *config
	mutex       *sync.RWMutex
	LatestQuote *StockQuote
	quoteChan   chan *StockQuote
	fillChan    chan *Execution
}

func InitSession(config *config, venue string) *Session {
	return &Session{
		Cash: 0,
		Position: 0,
		NAV: 0,
		Venue: venue,
		api: InitApi(config),
		config: config,
		mutex: &sync.RWMutex{},
		quoteChan: make(chan *StockQuote, 100),
		fillChan: make(chan *Execution, 100),
	}
}

func (o *Session) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("Session(%s)", string(j))
}

func (o *Session) Update(status *StockOrderAccountStatus) {
	var totalPosition int = 0
	var totalCash int = 0

	if !status.Ok {
		return
	}

	for _, so := range status.Orders {
		if !so.Ok {
			continue
		}

		for _, fill := range so.Fills {

			if so.Direction == DirectionBuy {
				totalCash -= fill.Price * fill.Qty
				totalPosition += fill.Qty
			}
			if so.Direction == DirectionSell {
				totalCash += fill.Price * fill.Qty
				totalPosition -= fill.Qty
			}
		}
	}

	o.mutex.Lock()
	if o.LatestQuote != nil && o.LatestQuote.Last > 0 {
		o.NAV = totalCash + (totalPosition * o.LatestQuote.Last)
	}
	o.Cash = totalCash
	o.Position = totalPosition
	o.mutex.Unlock()
}

func (o *Session) Observe(symbol string) {
	go o.api.StockTickerTape(o.quoteChan, o.Venue, symbol)
	go func(c chan *StockQuote) {
		for q := range c {
			o.mutex.Lock()
			o.LatestQuote = q
			o.mutex.Unlock()
		}
	}(o.quoteChan)

	go o.api.StockExecutions(o.fillChan, o.Venue, symbol)
	go func(c chan *Execution) {
		for e := range c {
			fmt.Println(e)
			if status, err := o.api.StockOrdersAccountStatus(o.Venue, symbol); err != nil {
				fmt.Println(err)
			} else {
				o.Update(status)
			}
		}
	}(o.fillChan)
}