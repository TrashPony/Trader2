package traderBot

import (
	"../traderInfo"
	"./Analyze"
	"./Worker"
)

const FEE = 0.0026 // в процентах 0.26, и сделано множителем тоесть поделено на 100

func Run(account *traderInfo.Account) {

	availableMarket := make([]*traderInfo.Market, 0)

	for _, marketBalance := range account.Balances {
		if !marketBalance.Available.IsZero() && marketBalance.Currency != "BTC" && marketBalance.Currency != "USDT" {
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

	if account.StartBTC >= 0.00075 {

		go traderInfo.UpdateActualMarketPool() // обновляет список всех маркетов, и обновляется каждые 2 часа

		baseInAlgorithm := &Analyze.AnalyzerInTrade{Name: "BaseInAlgorithm"}
		baseOutAlgorithm := &Analyze.AnalyzerOutTrade{Name: "BaseOutAlgorithm"}

		newWorker := &Worker.Worker{StartBTCCash: 0.00075, InTradeStrategy: baseInAlgorithm, OutTradeStrategy: baseOutAlgorithm}
		newWorker.Run(FEE)

	} else {
		println("Ты бомж :(")
	}
}
