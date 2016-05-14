package stockfighter

import (
	"github.com/franela/goreq"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
)

type Api struct {
	config *config
}

type HeartBeat struct {
	Ok    bool
	Error string
}

type VenueHeartBeat struct {
	Ok    bool
	Venue string
}

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

type Fill struct {
	Price int
	Qty   int
	Ts    time.Time
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

type StockOrderAccountStatus struct {
	Ok     bool
	Venue  string
	Orders []*StockOrder
}

func InitApi(config *config) *Api {
	return &Api{
		config: config,
	}
}

func (s *Api) HeartBeat() (*HeartBeat, error) {
	buffer, err := s.GetRequest("ob/api/heartbeat")
	if err != nil {
		return nil, err
	}

	var value *HeartBeat
	jsonErr := json.Unmarshal(buffer, &value)

	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) VenueHeartBeat(venue string) (*VenueHeartBeat, error) {
	buffer, err := s.GetRequest(fmt.Sprintf("ob/api/venues/%s/heartbeat", venue))
	if err != nil {
		return nil, err
	}

	var value *VenueHeartBeat

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) VenueStocks(venue string) (*VenueStocks, error) {
	buffer, err := s.GetRequest(fmt.Sprintf("ob/api/venues/%s/stocks", venue))
	if err != nil {
		return nil, err
	}

	var value *VenueStocks

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) StockOrderBook(venue string, stock string) (*StockOrderBook, error) {
	buffer, err := s.GetRequest(fmt.Sprintf("ob/api/venues/%s/stocks/%s", venue, stock))

	if err != nil {
		return nil, err
	}

	var value *StockOrderBook

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) StockOrder(soReq *StockOrderRequest) (*StockOrder, error) {
	url := fmt.Sprintf(
		"ob/api/venues/%s/stocks/%s/orders", soReq.Venue, soReq.Stock,
	)
	buffer, err := s.PostRequest(url, soReq)
	if err != nil {
		return nil, err
	}

	var value *StockOrder

	jsonErr := json.Unmarshal(buffer, &value)

	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) StockOrderCancel(so *StockOrder) (*StockOrder, error) {
	url := fmt.Sprintf(
		"ob/api/venues/%s/stocks/%s/orders/%d", so.Venue, so.Symbol, so.Id,
	)
	buffer, err := s.DeleteRequest(url)
	if err != nil {
		return nil, err
	}

	var value *StockOrder

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) StockOrdersAccountStatus(so *StockOrder) (*StockOrderAccountStatus, error) {
	urlFormat := "ob/api/venues/%s/accounts/%s/stocks/%s/orders"
	buffer, err := s.GetRequest(fmt.Sprintf(urlFormat, so.Venue, so.Account, so.Symbol))
	if err != nil {
		return nil, err
	}

	var value *StockOrderAccountStatus

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) StockQuote(venue string, stock string) (*StockQuote, error) {
	buffer, err := s.GetRequest(fmt.Sprintf("ob/api/venues/%s/stocks/%s/quote", venue, stock))
	if err != nil {
		return nil, err
	}

	var value *StockQuote

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) StockOrderStatus(so *StockOrder) (*StockOrder, error) {
	urlFormat := "ob/api/venues/%s/stocks/%s/orders/%d"
	buffer, err := s.GetRequest(fmt.Sprintf(urlFormat, so.Venue, so.Symbol, so.Id))
	if err != nil {
		return nil, err
	}

	var value *StockOrder

	jsonErr := json.Unmarshal(buffer, &value)
	if jsonErr == nil {
		return value, nil
	}

	return nil, jsonErr
}

func (s *Api) GetRequest(url string) ([]byte, error) {
	return s.Request("GET", url, nil)
}

func (s *Api) DeleteRequest(url string) ([]byte, error) {
	return s.Request("DELETE", url, nil)
}

func (s *Api) PostRequest(url string, body interface{}) ([]byte, error) {
	return s.Request("POST", url, body)
}

func (s *Api) Request(method string, path string, body interface{}) ([]byte, error) {

	if body == nil {
		body = ""
	}

	req := goreq.Request{
		Accept: "application/json",
		Body: body,
		ContentType: "application/json",
		Method: method,
		Uri: fmt.Sprintf("https://api.stockfighter.io/%s", path),

	}
	req.AddHeader("X-Starfighter-Authorization", s.config.ApiKey)

	response, err := req.Do()
	if err != nil {
		return nil, err
	}

	buffer, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	} else {
		fmt.Println(string(buffer))
	}
	return buffer, nil
}