package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pavben/elise/utils/csv"
)

import (
	"./controllers"
)

func main() {

	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := pool.Get()

	values, err := csv.ReadCsvFile("./venders/serviceconfig/marketList.csv")
	fmt.Print(err)
	marketArray := []agent.Market{}
	fmt.Print("superwolf")
	fmt.Print(values)
	for i := 0; i < len(values); i++ {
		rate, _ := strconv.ParseFloat(values[i][4], 64)
		market := agent.Market{
			Id:    values[i][0],
			Name:  values[i][1],
			Code:  values[i][2],
			Agent: values[i][3],
			Rate:  rate,
		}
		marketArray = append(marketArray, market)
		jsonString, _ := json.Marshal(market)
		fmt.Println(string(jsonString))
		_, err := conn.Do("hset", "superwolf", market.Id, jsonString)

		fmt.Print(err)
	}

}
