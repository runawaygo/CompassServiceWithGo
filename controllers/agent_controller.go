package agent

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strings"
)

type Market struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Code  string  `json:"code"`
	Agent string  `json:"agent"`
	Rate  float64 `json:"rate"`
}

type MarketRepo interface {
	GetMarketByMarketId(marketId string) Market
	SetMarketAgentByMarketId(marketId string, agent string) Market
}

type MarketMap map[string]Market

func (m MarketMap) GetMarketByMarketId(marketId string) Market {
	return m[marketId]
}

func (m MarketMap) SetMarketAgentByMarketId(marketId string, agent string) Market {
	market := m[marketId]
	market.Agent = agent

	m[marketId] = market

	return market
}

func newMarket() *Market {
	market := new(Market)
	market.Rate = 0
	return market
}

func SetMarket(request *http.Request, conn redis.Conn, marketRepo MarketRepo) string {
	id, agent := request.FormValue("id"), request.FormValue("agent")
	marketRepo.SetMarketAgentByMarketId(id, agent)

	market := marketRepo.GetMarketByMarketId(id)
	marketJson, _ := json.Marshal(market)

	return string(marketJson)
}

func GetMarket(request *http.Request, conn redis.Conn) string {
	var market Market
	qs := request.URL.Query()
	marketId := qs.Get("marketId")
	reply, _ := redis.Bytes(conn.Do("get", marketId))
	json.Unmarshal(reply, &market)
	return string(reply)
}

func GetAllMarket(conn redis.Conn) string {
	values, _ := redis.Strings(conn.Do("hgetall", "superwolf"))
	marketArray := []string{}
	for i := 1; i < len(values); i += 2 {
		fmt.Println(values[i])
		marketArray = append(marketArray, values[i])
	}

	resultString := "[" + strings.Join(marketArray, ",") + "]"
	return resultString
}

func Agent(request *http.Request, marketRepo MarketRepo) string {
	qs := request.URL.Query()
	marketId := qs.Get("marketId")
	if marketId == "122" {
		panic("superwolf")
	}
	market := marketRepo.GetMarketByMarketId(marketId)
	marketJson, _ := json.Marshal(market)
	return string(marketJson)
}

func Agent2(request *http.Request, conn redis.Conn) string {
	var market Market
	qs := request.URL.Query()
	marketId := qs.Get("marketId")
	reply, _ := redis.Bytes(conn.Do("get", marketId))

	json.Unmarshal(reply, &market)
	return string(reply)
}
