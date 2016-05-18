package stockfighter

import (
	"github.com/franela/goreq"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/websocket"
)

/*
Web Sockets: https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/08.2.html
*/

type Api struct {
	Config  *config
}

func InitApi(config *config) *Api {
	return &Api{
		Config: config,
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
	req.AddHeader("X-Starfighter-Authorization", s.Config.ApiKey)

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

func (s *Api) VenueTickerTape(stockQuoteChan chan *StockQuote, venue string) error {
	urlFormat := "ob/api/ws/%s/venues/%s/tickertape"
	url := fmt.Sprintf(urlFormat, s.Config.Account, venue)
	return s.wsStockQuote(stockQuoteChan, url)
}

func (s *Api) StockTickerTape(stockQuoteChan chan *StockQuote, venue string, stock string) error {
	urlFormat := "ob/api/ws/%s/venues/%s/tickertape/stocks/%s"
	url := fmt.Sprintf(urlFormat, s.Config.Account, venue, stock)
	return s.wsStockQuote(stockQuoteChan, url)
}

func (s *Api) wsStockQuote(stockQuoteChan chan *StockQuote, url string) error {
	conn, err := s.Stream(url)
	if err != nil {
		return err
	}
	for {
		var sStockQuote *wrappedStockQuote
		if err := websocket.JSON.Receive(conn, &sStockQuote); err != nil {
			fmt.Println("message error:", err)
			continue
		}
		//fmt.Printf("Received StreamQuote: %#v\n", sStockQuote)
		stockQuoteChan <- sStockQuote.Quote
	}
	return nil
}

func (s *Api) VenueExecutions(executionsChan chan *Execution, venue string) error {
	urlFormat := "ob/api/ws/%s/venues/%s/executions"
	url := fmt.Sprintf(urlFormat, s.Config.Account, venue)
	return s.wsExecutions(executionsChan, url)
}

func (s *Api) StockExecutions(executionsChan chan *Execution, venue string, stock string) error {
	urlFormat := "ob/api/ws/%s/venues/%s/executions/stocks/%s"
	url := fmt.Sprintf(urlFormat, s.Config.Account, venue, stock)
	return s.wsExecutions(executionsChan, url)
}

func (s *Api) wsExecutions(executionsChan chan *Execution, url string) error {
	conn, err := s.Stream(url)
	if err != nil {
		return err
	}
	for {
		var execution *Execution
		if err := websocket.JSON.Receive(conn, &execution); err != nil {
			fmt.Println("message error:", err)
			continue
		}
		//fmt.Printf("Received Execution: %#v\n", execution)
		executionsChan <- execution
	}
	return nil
}

func (s *Api) Stream(path string) (ws *websocket.Conn, err error) {
	var origin = "https://api.stockfighter.io/"
	var url = fmt.Sprintf("wss://api.stockfighter.io/%s", path)
	return websocket.Dial(url, "", origin)
}