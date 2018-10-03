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
				priceByu, priceByuOk := market.OrdersBuy[0].Rate.Float64()

				if ok && volume > 20 && AnalyzerInTrade(market) && btcOk && priceByuOk {
					fee := availableBTC * 0.0026
					quantity := (availableBTC - fee) / priceByu

					uuidBuyOrder, err = market.BuyLimit(decimal.NewFromFloat(quantity), decimal.NewFromFloat(priceByu))
					if err != nil {
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

			}

			for _, order := range orders {
				if order.OrderUuid == uuidBuyOrder { // если находим заказ, значит не купили и снимаем его
					err = buyMarket.CancelOrder(uuidBuyOrder)
					if err != nil {

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
