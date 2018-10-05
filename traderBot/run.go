package traderBot

import (
	"../traderInfo"
	"./Worker"
	"./Analyze"
)

const FEE = 0.0026 // в процентах 0.26, и сделано множителем тоесть поделено на 100

func Run(account *traderInfo.Account) {

	availableMarket := make([]*traderInfo.Market, 0)

	for _, marketBalance := range account.Balances {
		if !marketBalance.Available.IsZero() && marketBalance.Currency != "BTC" && marketBalance.Currency != "USDT" {
			availableMarket = append(availableMarket, traderInfo.GetMarket("BTC-"+marketBalance.Currency))
		} else {
			var ok bool
			account.StartBTC, ok = marketBalance.Available.Float64()
			if ok && marketBalance.Currency == "BTC" {
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

	if account.StartBTC > 0.0005 {
		// пока только 1 бот и у него нет денег :)

		baseInAlgorithm := &Analyze.AnalyzerInTrade{Name: "BaseInAlgorithm"}
		baseOutAlgorithm := &Analyze.AnalyzerOutTrade{Name: "BaseOutAlgorithm"}

		newWorker := &Worker.Worker{StartBTCCash: account.StartBTC, InTradeStrategy: baseInAlgorithm, OutTradeStrategy: baseOutAlgorithm}
		newWorker.Run(FEE)
	} else {
		println("Ты бомж :(")
	}
}
