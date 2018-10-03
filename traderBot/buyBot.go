package traderBot

import "../traderInfo"

func TradeBuyBot(account *traderInfo.Account) {
	for {
		markets := traderInfo.GetAllMarket()

		if markets == nil {
			continue
		}

		for _, market := range markets {
			volume, ok := market.MarketSummary.Volume.Float64()

			if ok && volume > 20 && AnalyzerInTrade(market) {
				// заходим на рынок
			}
		}
	}
}
