package traderBot

import "../traderInfo"

import "time"

func TradeSellBot(market *traderInfo.Market, account *traderInfo.Account) {

	var uuidSellOrder string

	for {
		var err error

		market.UpdateMarket()   // каждую итерацию обновляем рынок
		account.UpdateAccount() // и аканут

		if AnalyzerOutTrade(market) { // провека на необходимость выхода из рынка
			// тут пыатется сбыть по цене 1го ордера покупки, это надо делать если надо очень быстро выйти из трейда
			// TODO в идиале надо пытаться продавать чуть дешевле чем первый ордер на продажу
			uuidSellOrder, err = market.SellLimit(account.GetAvailableCurrencyBalance(market.CurrencyPair), &market.OrdersBuy[0].Rate)
			if err != nil {
				continue // TODO
			}
		} else {
			if uuidSellOrder != "" {
				err = market.CancelOrder(uuidSellOrder)
				if err != nil {
					// TODO
				} else {
					uuidSellOrder = ""
				}
			}
		}

		if account.GetAvailableCurrencyBalance(market.CurrencyPair).IsZero() {
			return // если монет не осталось выходим из трейда
		}

		time.Sleep(200 * time.Millisecond) // на всякий случай задержка, а то забанят еще)
	}
}
