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
}