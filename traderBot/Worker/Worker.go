package Worker

import (
	"../../traderInfo"
	"../Analyze"
	"strconv"
	"sync"
)

type Worker struct {
	InTradeStrategy   *Analyze.AnalyzerInTrade      `json:"in_trade_strategy"`
	OutTradeStrategy  *Analyze.AnalyzerOutTrade     `json:"out_trade_strategy"`
	StartBTCCash      float64                       `json:"start_btc_cash"`
	AvailableBTCCash  float64                       `json:"available__btc_cash"`
	BuyActiveMarket   *traderInfo.Market            `json:"active_markets"`
	SellActiveMarkets map[string]*traderInfo.Market `json:"sell_active_markets"`
	AltBalances       map[string]*Alt               `json:"alt_balances"`
	Fee               float64                       `json:"fee"`
	mx                sync.Mutex
}

func (worker *Worker) Run() bool {
	// получить кеш бота, принять анализаторы, сделать проверки на возможность рынка и запуск горутины работы
	// добавлять их в
	return true
}

type Alt struct {
	AltName  string  `json:"alt_name"`
	Balance  float64 `json:"balance"`
	BuyPrice float64 `json:"buy_price"` // цена за которую купил эту пачку
}

func (worker *Worker) AddAlt(altName string, balance, buyPrice float64) {

	worker.mx.Lock()
	defer worker.mx.Unlock()

	alt := &Alt{AltName: altName, Balance: balance, BuyPrice: buyPrice}
	// немного ебаный ключь мапы ¯\_(ツ)_/¯
	worker.AltBalances[altName+":"+strconv.FormatFloat(balance, 'f', 6, 64)+":"+strconv.FormatFloat(buyPrice, 'f', 6, 64)] = alt
}

func (worker *Worker) RemoveAlt(alt *Alt) {

	worker.mx.Lock()
	defer worker.mx.Unlock()

	key := alt.AltName + ":" + strconv.FormatFloat(alt.Balance, 'f', 6, 64) + ":" + strconv.FormatFloat(alt.BuyPrice, 'f', 6, 64)
	delete(worker.AltBalances, key)
}