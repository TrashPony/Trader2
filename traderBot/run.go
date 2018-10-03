package traderBot

import (
	"../traderInfo"
)

func Run(account *traderInfo.Account) {

	availableMarket := make([]*traderInfo.Market, 0)

	for _, marketBalance := range account.Balances {
		if !marketBalance.Available.IsZero() && marketBalance.Currency != "BTC" && marketBalance.Currency != "USDT" {
			availableMarket = append(availableMarket, traderInfo.GetMarket("BTC-"+marketBalance.Currency))
		} else {
			available, ok := marketBalance.Available.Float64()
			if ok && marketBalance.Currency == "BTC" {
				print("available BTC: ")
				println(available)
			}
		}
	}

	if len(availableMarket) > 0 { // если уже приобретены валюты то надо начать торговать с них, в идиале их быть не должно, всегда выходить в чистый BTC надо
		for _, market := range availableMarket {
			if market != nil {
				go TradeSellBot(market, account, market.OrdersSell[0].Rate) // запуска торгового бота на продаху, будет работать в отдельном потоке и не мешать основному боту
			}
		}
	}
}