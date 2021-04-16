package traderBot

import (
	"github.com/TrashPony/Trader2/traderBot/Analyze"
	"github.com/TrashPony/Trader2/traderBot/Worker"
	"github.com/TrashPony/Trader2/traderInfo"
)

const FEE = 0.0075 // в процентах 0.75, и сделано множителем тоесть поделено на 100
const minBTC = 0.0006

func Run(account *traderInfo.Account) {

	availableMarket := make([]*traderInfo.Market, 0)

	for _, marketBalance := range account.Balances {
		if !marketBalance.Available.IsZero() && marketBalance.Currency != "BTC" && marketBalance.Currency != "USDT" && marketBalance.Currency != "USD" {
			availableMarket = append(availableMarket, traderInfo.GetMarket("BTC-"+marketBalance.Currency))
		} else {
			if marketBalance.Currency == "BTC" {
				account.StartBTC, _ = marketBalance.Available.Float64()
				print("available BTC: ")
				println(account.StartBTC)
			}
		}
	}

	if len(availableMarket) > 0 { // если уже приобретены валюты то надо начать торговать с них, в идиале их быть не должно, всегда выходить в чистый BTC надо
		for _, market := range availableMarket {
			if market != nil {
				// TODO сливать старые монеты по рынку
			}
		}
	}

	if account.StartBTC >= minBTC { // * 3

		go traderInfo.UpdateActualMarketPool() // обновляет список всех маркетов, и обновляется каждые 2 часа

		fastInAlgorithm := &Analyze.AnalyzerInTrade{Name: "fast"}
		fastOutAlgorithm := &Analyze.AnalyzerOutTrade{Name: "fast"}

		newWorker := &Worker.Worker{TradeStrategy: "Fast", StartBTCCash: minBTC + 0.00001, InTradeStrategy: fastInAlgorithm, OutTradeStrategy: fastOutAlgorithm}
		go newWorker.Run(FEE)

		newWorker2 := &Worker.Worker{TradeStrategy: "Slow", StartBTCCash: minBTC + 0.00001, InTradeStrategy: fastInAlgorithm, OutTradeStrategy: fastOutAlgorithm}
		go newWorker2.Run(FEE)

		newWorker3 := &Worker.Worker{TradeStrategy: "VerySlow", StartBTCCash: minBTC + 0.00001, InTradeStrategy: fastInAlgorithm, OutTradeStrategy: fastOutAlgorithm}
		go newWorker3.Run(FEE)

	} else {
		println("Ты бомж :(")
	}
}
