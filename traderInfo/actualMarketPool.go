package traderInfo

import (
	"strings"
	"sync"
	"time"
)

type MarketPool struct {
	Markets map[string]*Market
	mx      sync.Mutex
}

var Markets = NewMarketPool()

func NewMarketPool() *MarketPool {
	return &MarketPool{
		Markets: make(map[string]*Market),
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

func UpdateActualMarketPool() {
	for {
		markets, err := GetBittrex().GetMarkets()
		if err != nil {
			println(err.Error())
		}

		for _, market := range markets {

			pair := strings.Split(market.MarketName, "-")

			if pair[0] == "BTC" && market.IsActive {

				newMarket := GetMarket(market.MarketName)

				volume, _ := newMarket.MarketSummary.BaseVolume.Float64()

				if volume > 20 { // если у валюты капитализация больше 20 биткоинов
					Markets.AddMarket(newMarket)
				}
			}
		}
		time.Sleep(time.Hour * 2) // каждые 2 часа полностью обновляет список достпуынйх маркетов
	}
}
