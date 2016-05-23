package stockfighter

import (
	"time"
	"fmt"
	"encoding/json"
)

const (
	DirectionBuy string = "buy"
	DirectionSell string = "sell"
	OrderLimit string = "limit"
	OrderMarket string = "market"
	OrderFillOrKill string = "fill-or-kill"
	OrderImmediateOrCancel string = "immediate-or-cancel"
)

type Symbol struct {
	Name   string
	Symbol string
}

type Bid struct {
	Price int
	Qty   int
	IsBuy bool
}

type Ask struct {
	Price int
	Qty   int
	IsBuy bool
}

type Fill struct {
	Price int
	Qty   int
	Ts    time.Time
}

type HeartBeat struct {
	Ok    bool
	Error string
}

type VenueHeartBeat struct {
	Ok    bool
	Venue string
}

type VenueStocks struct {
	Ok      bool
	Symbols []*Symbol
}

type StockOrderBook struct {
	Asks   []*Ask
	Bids   []*Bid
	Ok     bool
	Symbol string
	Ts     time.Time
	Venue  string
}

type StockOrder struct {
	Account     string
	Direction   string
	Fills       []*Fill
	Id          int
	Ok          bool
	Open        bool
	OrderType   string
	OriginalQty int
	Price       int
	Qty         int
	Symbol      string
	TotalFilled int
	Ts          time.Time
	Venue       string
}

type StockQuote struct {
	Ok        bool
	Symbol    string
	Venue     string
	Bid       int
	Ask       int
	BidSize   int
	AskSize   int
	BidDepth  int
	AskDepth  int
	Last      int
	LastSize  int
	LastTrade time.Time
	QuoteTime time.Time
}

type Execution struct {
	Ok               bool
	Account          string
	Venue            string
	Symbol           string
	Order            *StockOrder
	StandingId       int
	IncomingId       int
	Price            int
	Filled           int
	FilledAt         time.Time
	StandingComplete bool
	IncomingComplete bool
}

type wrappedStockQuote struct {
	Ok    bool
	Quote *StockQuote
}

type StockOrderRequest struct {
	Account   string
	Direction string
	OrderType string
	Price     int
	Qty       int
	Stock     string
	Venue     string
}

type StockOrderAccountStatus struct {
	Ok     bool
	Venue  string
	Orders []*StockOrder
}

func (o *StockOrderAccountStatus) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("StockOrderAccountStatus(%s)", string(j))
}

func (o *StockOrderRequest) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("StockOrderRequest(%s)", string(j))
}

func (o *Execution) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("Execution(%s)", string(j))
}

func (o *StockOrder) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("StockOrder(%s)", string(j))
}

func (o *StockQuote) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("StockQuote(%s)", string(j))
}

func (o *Fill) String() string {
	j, _ := json.Marshal(o)
	return fmt.Sprintf("Fill(%s)", string(j))
}

