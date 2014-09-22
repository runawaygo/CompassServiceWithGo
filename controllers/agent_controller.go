package agent

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	// "github.com/pavben/elise/utils/csv"
	"net/http"
)

type Market struct {
	Id    string  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Code  string  `json:"code,omitempty"`
	Agent string  `json:"agent,omitempty"`
	Rate  float64 `json:"rate,omitempty"`
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

func GetMockMarket() *Market {
	market := newMarket()
	market.Id = "123"
	market.Name = "superwolf"
	market.Code = "superowlf"
	market.Agent = "superwolf"
	return market
}

func SetAgent(conn redis.Conn, marketRepo MarketRepo) string {
	fmt.Println(marketRepo.GetMarketByMarketId("12"))
	marketRepo.SetMarketAgentByMarketId("12", "tt")
	return ""
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
