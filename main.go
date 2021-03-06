package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/pavben/elise/utils/csv"
)

import (
	"./controllers"
)

func main() {
	m := martini.Classic()
	m.Use(redismw("tcp", "127.0.0.1:6379"))
	m.Use(agentMarketMap())

	m.Get("/r ", func(request *http.Request, conn redis.Conn) string {
		qs := request.URL.Query()
		conn.Do("set", "name", "abc")
		reply, err := redis.String(conn.Do("get", "name"))
		if err != nil {
			return "error"
		}
		return reply + qs.Get("name")
	})
	m.Get("/preload/config", func(res http.ResponseWriter) {
		out, err := os.Open("./venders/serviceconfig/preload.config.json")

		if err != nil {
			log.Fatalln(err)
		}
		io.Copy(res, out)
	})
	m.Get("/agent", agent.GetMarket)
	m.Post("/agent", agent.SetMarket)
	m.Get("/agent/index", agent.GetAllMarket)
	m.Run()
}

func agentMarketMap() martini.Handler {
	marketMap := agent.MarketMap{}
	values, _ := csv.ReadCsvFile("./venders/serviceconfig/marketList.csv")
	for i := 0; i < len(values); i++ {
		rate, _ := strconv.ParseFloat(values[i][4], 64)
		marketMap[values[i][0]] = agent.Market{
			Id:    values[i][0],
			Name:  values[i][1],
			Code:  values[i][2],
			Agent: values[i][3],
			Rate:  rate,
		}
	}
	return func(c martini.Context) {
		c.MapTo(&marketMap, (*agent.MarketRepo)(nil))
	}
}

func redismw(proto, addr string) martini.Handler {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(proto, addr)
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
	_, err := pool.Get().Do("PING")
	if err != nil {
		panic(err.Error())
	}
	println("Init Redis middleware successfully.")
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		c.MapTo(pool.Get(), (*redis.Conn)(nil))
	}
}
