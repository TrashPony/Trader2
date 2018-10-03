package traderBot

import (
	"../traderInfo"
	"github.com/shopspring/decimal"
	"time"
)

func TradeBuyBot(account *traderInfo.Account) {

	uuidBuyOrder := ""
	var buyMarket *traderInfo.Market
	var buyRate float64

	for {

		var err error

		if uuidBuyOrder == "" && buyMarket == nil {
			markets := traderInfo.GetAllMarket()

			if markets == nil {
				continue
			}

			for _, market := range markets {

				volume, ok := market.MarketSummary.Volume.Float64()
				availableBTC, btcOk := account.GetAvailableCurrencyBalance("BTC").Float64()

				buy, fast := AnalyzerInTrade(market)
				var priceByu float64
				var priceByuOk bool

				if fast { // если фаст значит в рынок надо войти как можно быстрее )
					priceByu, priceByuOk = market.OrdersSell[0].Rate.Float64()
				} else {
					priceByu, priceByuOk = market.OrdersBuy[0].Rate.Float64()
					priceByu += 0.00000001 // наращиваем 1 сатоши что бы стать самым первым ордером в стакане
				}

				if ok && volume > 20 && buy && btcOk && priceByuOk {
					fee := availableBTC * FEE
					quantity := (availableBTC - fee) / priceByu

					uuidBuyOrder, err = market.BuyLimit(decimal.NewFromFloat(quantity), decimal.NewFromFloat(priceByu))
					if err != nil {
						println(err.Error())
						continue
					}

					buyMarket = market
					buyRate = priceByu
					break
				}
			}
		} else {

			time.Sleep(time.Second * 5) // задержка что бы купили валюту

			orders, err := buyMarket.GetOpenOrders()
			if err != nil {
				println(err.Error())
			}

			for _, order := range orders {
				if order.OrderUuid == uuidBuyOrder { // если находим заказ, значит не купили и снимаем его
					err = buyMarket.CancelOrder(uuidBuyOrder)
					if err != nil {
						println(err.Error())
					} else {
						uuidBuyOrder = ""
						buyMarket = nil
						buyRate = 0
						continue
					}
				}
			}

			// если мы дошли до сюда значит монеты у нас купили и их надо продать :)
			go TradeSellBot(buyMarket, account, decimal.NewFromFloat(buyRate))
			break
		}
	}
}
