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
		println("e")

		market.UpdateMarket()   // каждую итерацию обновляем рынок
		account.UpdateAccount() // и аканут

		sell, fast := AnalyzerOutTrade(market)

		if sell {
			// TODO проверки на цену покупки и профита
			// TODO проверка на комисию
			if fast { // если фаст значит из рынка нужно выйти немедленно
				uuidSellOrder, err = market.SellLimit(account.GetAvailableCurrencyBalance(strings.Split(market.CurrencyPair, "-")[1]), &market.OrdersBuy[0].Rate)
				if err != nil {
					println(err.Error())
				}
			} else {
				// TODO
			}
		} else {
			if uuidSellOrder != "" {
				err = market.CancelOrder(uuidSellOrder)
				if err != nil {
					println(err.Error())
				} else {
					uuidSellOrder = ""
				}
			}
		}

		if account.GetAvailableCurrencyBalance(strings.Split(market.CurrencyPair, "-")[1]).IsZero() {
			return // если монет не осталось выходим из трейда
		}
	}
}
