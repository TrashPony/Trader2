package traderInfo

import (
	"github.com/shopspring/decimal"
	"strings"
	"sync"
	"time"
)

type MarketPool struct {
	Markets      map[string]*Market
	minTradeSize map[string]decimal.Decimal
	mx           sync.Mutex
}

var Markets = NewMarketPool()

func NewMarketPool() *MarketPool {
	return &MarketPool{
		Markets:      make(map[string]*Market),
		minTradeSize: make(map[string]decimal.Decimal),
	}
}

func (marketPoll *MarketPool) GetMarketPool() map[string]*Market {

	marketPoll.mx.Lock()
	defer marketPoll.mx.Unlock()

	return marketPoll.Markets
}

func (marketPoll *MarketPool) AddMarket(newMarket *Market) {
	marketPoll.mx.Lock()
	defer marketPoll.mx.Unlock()
	marketPoll.Markets[newMarket.CurrencyPair] = newMarket
}

func (marketPoll *MarketPool) GetMinTradeSize(mName string) float64 {
	market, ok := Markets.minTradeSize[mName]
	if ok {
		minTradeFloat, _ := market.Float64()
		return minTradeFloat
	}

	return 0
}

func UpdateActualMarketPool() {

	markets, _ := GetBittrex().GetMarkets()

	for _, market := range markets {
		Markets.minTradeSize[market.MarketName] = market.MinTradeSize
	}

	for {
		markets, err := GetBittrex().GetMarketSummaries()
		if err != nil {
			println(err.Error())
		}

		for _, marketSummaries := range markets {

			pair := strings.Split(marketSummaries.MarketName, "-")

			if pair[0] == "BTC" && pair[1] != "USDT" && pair[1] != "USD" {
				volume, _ := marketSummaries.BaseVolume.Float64()
				if volume >= 5 { // если у валюты капитализация больше 5 биткоинов
					newMarket := GetMarket(marketSummaries.MarketName)
					if newMarket != nil && len(newMarket.OrdersBuy) > 0 {
						Markets.AddMarket(newMarket)
					}
				}
			}
		}
		time.Sleep(time.Hour * 2) // каждые 2 часа полностью обновляет список достпуынйх маркетов
	}
}
