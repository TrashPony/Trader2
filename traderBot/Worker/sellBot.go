package Worker

import (
	"../Analyze"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func (worker *Worker) TradeSellBot() {

	var uuidSellOrder string
	var sellRate float64

	for {
		for _, market := range worker.SellActiveMarkets { // крутим все активные маркеты где совершили покупку
			for _, altBalances := range worker.AltBalances { // берем прошлые сделки покупок
				if altBalances.AltName == strings.Split(market.CurrencyPair, "-")[1] { // проверяем что алт баланс соотвевтует маркету
					var err error

					market.UpdateMarket()

					sell, fast, newProfit := worker.OutTradeStrategy.Analyze(market, altBalances.ProfitPrice, altBalances.GrowProfitPrice)

					altBalances.GrowProfitPrice = newProfit

					if sell {

						var priceSell float64
						var priceSellOk bool

						if fast { // если фаст значит из рынка нужно выйти немедленно
							priceSell, priceSellOk = market.OrdersBuy[0].Rate.Float64()
						} else {
							priceSell, priceSellOk = market.OrdersSell[0].Rate.Float64()
							if priceSellOk {
								priceSell -= 0.00000001 // уменьшаем на 1 сатоши что бы стать самым первым ордером в стакане
							}
						}

						if priceSellOk {
							uuidSellOrder, err = market.SellLimit(decimal.NewFromFloat(altBalances.Balance), decimal.NewFromFloat(priceSell))
							if err != nil {
								println(err.Error())
							}

							sellRate = priceSell
						}

					} else {

						market.UpdateMarket()

						orders, err := market.GetOpenOrders()
						if err != nil {
							println(err.Error())
						}

						findOrder := false

						for _, order := range orders {
							if order.OrderUuid == uuidSellOrder { // если находим заказ, значит не купили или купили частично
								findOrder = true

								orderQuantity, _ := order.Quantity.Float64()

								if orderQuantity == altBalances.Balance {
									// не выкупили

									first, err := market.GetFirstSellOrder()
									if err != nil {
									} else {

										firstRate, okRate := first.Rate.Float64()
										sell, _, newProfit := worker.OutTradeStrategy.Analyze(market, altBalances.ProfitPrice, altBalances.GrowProfitPrice)
										altBalances.GrowProfitPrice = newProfit

										if okRate && firstRate < sellRate && sell {
											// если у первого заказа цена выше чем у наc
											// и рынок досихпор считается перспективным то пересоздаем ордер
											err = market.CancelOrder(uuidSellOrder)
											if err != nil {
												println(err.Error())
											} else {
												uuidSellOrder = ""
												sellRate = 0
											}
										}
									}
								} else {
									// частичный выкуп
									buyAltCount := altBalances.Balance - orderQuantity // высчитываем сколько купили
									worker.AvailableBTCCash += buyAltCount * sellRate  // высчитывает заработаные BTC
									altBalances.Balance -= buyAltCount                 // отнимать у баланса валюты проданые монеты

									fmt.Print("Продалась часть закупа с выгодой ", Analyze.PercentageCalculator(altBalances.ProfitPrice, sellRate))
								}
							}
						}

						if !findOrder {
							// выкупили полностью
							worker.AvailableBTCCash += altBalances.Balance * sellRate
							uuidSellOrder = ""
							sellRate = 0

							fmt.Print("Продалось относительно начального закупа с выгодой ", Analyze.PercentageCalculator(altBalances.ProfitPrice, sellRate))

							worker.RemoveAlt(altBalances)
						}
					}
				}
			}
		}
	}
}
