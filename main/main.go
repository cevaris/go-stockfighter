package main

import (
	"time"
	"fmt"
	"github.com/cevaris/stockfighter"
)

func main() {
	fmt.Println(time.Now().Unix())

	config := stockfighter.InitConfig(".env.yml")

	fmt.Printf("config: %#v\n", config.ApiKey)

	api := stockfighter.InitApi(config)

	if value, err := api.HeartBeat(); err == nil {
		fmt.Printf("request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	if value, err := api.VenueHeartBeat("TESTEX"); err == nil {
		fmt.Printf("request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	if value, err := api.VenueStocks("TESTEX"); err == nil {
		fmt.Printf("request: %#v\n", value)
		fmt.Printf("%#v\n", value.Symbols[0].Name)
	} else {
		fmt.Println(err)
	}

	if value, err := api.StockOrderBook("TESTEX", "FOOBAR"); err == nil {
		fmt.Printf("request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	soReq := &stockfighter.StockOrderRequest{
		Account: config.Account,
		Venue: "CEMEX",
		Stock: "ZICO",
		Price: 10,
		Qty: 2,
		Direction: "buy",
		OrderType: "limit",
	}
	stockOrder, soResErr := api.StockOrder(soReq);
	if soResErr == nil {
		fmt.Printf("stockorder request: %#v\n", stockOrder)
	} else {
		fmt.Println(soResErr)
	}

	if value, err := api.StockOrderCancel(stockOrder); err == nil {
		fmt.Printf("stock order cancel request: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	if value, err := api.StockOrderStatus(stockOrder); err == nil {
		fmt.Printf("stock order status: %#v\n", value)
	} else {
		fmt.Println(err)
	}

	if value, err := api.StockOrdersAccountStatus(stockOrder); err == nil {
		fmt.Printf("all stock order account status: %#v\n", value)
		for i, order := range value.Orders {
			fmt.Printf("order %d: %#v\n", i, order)
		}
	} else {
		fmt.Println(err)
	}

}