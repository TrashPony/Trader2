package traderBot

import (
	"../traderInfo"
	"github.com/shopspring/decimal"
	"strings"
)

func TradeSellBot(market *traderInfo.Market, account *traderInfo.Account, priceBuy decimal.Decimal) {

	var uuidSellOrder string

	for {
		var err error

		market.UpdateMarket()   // каждую итерацию обновляем рынок
		account.UpdateAccount() // и аканут

		if AnalyzerOutTrade(market) { // провека на необходимость выхода из рынка
			// тут пыатется сбыть по цене 1го ордера покупки, это надо делать если надо очень быстро выйти из трейда
			// TODO в идиале надо пытаться продавать чуть дешевле чем первый ордер на продажу
			// TODO проверки на цену покупки и профита
			uuidSellOrder, err = market.SellLimit(account.GetAvailableCurrencyBalance(strings.Split(market.CurrencyPair, "-")[1]), &market.OrdersBuy[0].Rate)
			if err != nil {

			}
		} else {
			if uuidSellOrder != "" {
				err = market.CancelOrder(uuidSellOrder)
				if err != nil {

				} else {
					uuidSellOrder = ""
				}
			}
		}

		if account.GetAvailableCurrencyBalance(market.CurrencyPair).IsZero() {
			return // если монет не осталось выходим из трейда
		}
	}
}
