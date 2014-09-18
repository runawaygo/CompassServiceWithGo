package main

import (
	// "github.com/paulchiu/gone-lib/csv"
	"./controllers"
	"fmt"
	"github.com/pavben/elise/utils/csv"
)

func main() {
	values, err := csv.ReadCsvFile("./verders/agent.csv")
	market := agent.GetMockMarket()
	// println(agent.Market.Elem())
	println(values)
	fmt.Print(err)
	fmt.Println("superwolf")
	fmt.Print(market)
}
