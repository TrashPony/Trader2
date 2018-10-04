package Worker

import (
	"../../traderInfo"
	"github.com/shopspring/decimal"
	"strings"
)

func (worker *Worker) TradeBuyBot() {

	uuidBuyOrder := ""
	var buyRate float64
	var buyQuantityAlt float64

	for {

		if worker.AvailableBTCCash >= 0.00050000 {

			var err error

			if uuidBuyOrder == "" && worker.BuyActiveMarket == nil {
				markets := traderInfo.GetAllMarket()

				if markets == nil {
					continue
				}

				for _, market := range markets {

					volume, ok := market.MarketSummary.Volume.Float64()
					buy, fast := worker.InTradeStrategy.Analyze(market)

					if ok && volume > 20 && buy {

						var priceByu float64
						var priceByuOk bool

						if fast { // если фаст значит в рынок надо войти как можно быстрее )
							priceByu, priceByuOk = market.OrdersSell[0].Rate.Float64()
						} else {
							priceByu, priceByuOk = market.OrdersBuy[0].Rate.Float64()
							if priceByuOk {
								priceByu += 0.00000001 // наращиваем 1 сатоши что бы стать самым первым ордером в стакане
							}
						}

						if priceByuOk {
							fee := worker.AvailableBTCCash * worker.Fee
							quantity := (worker.AvailableBTCCash - fee) / priceByu

							uuidBuyOrder, err = market.BuyLimit(decimal.NewFromFloat(quantity), decimal.NewFromFloat(priceByu))
							if err != nil {
								println(err.Error())
								continue
							}

							worker.BuyActiveMarket = market
							buyRate = priceByu
							buyQuantityAlt = quantity

							break
						}
					}
				}
			} else {

				worker.BuyActiveMarket.UpdateMarket()

				orders, err := worker.BuyActiveMarket.GetOpenOrders()
				if err != nil {
					println(err.Error())
				}

				findOrder := false

				for _, order := range orders {
					if order.OrderUuid == uuidBuyOrder { // если находим заказ, значит не купили или купили частично
						findOrder = true

						orderQuantity, _ := order.Quantity.Float64()

						if orderQuantity == buyQuantityAlt {
							// не выкупили

							first, err := worker.BuyActiveMarket.GetFirstBuyOrder()
							if err != nil {
							} else {

								firstRate, okRate :=  first.Rate.Float64()
								buy, _ := worker.InTradeStrategy.Analyze(worker.BuyActiveMarket)

								if okRate && firstRate > buyRate && buy {
									// если у первого заказа цена выше чем у наc
									// и рынок досихпор считается перспективным то пересоздаем ордер
									err = worker.BuyActiveMarket.CancelOrder(uuidBuyOrder)
									if err != nil {
										println(err.Error())
									} else {
										uuidBuyOrder = ""
										worker.BuyActiveMarket = nil
										buyRate = 0
										buyQuantityAlt = 0
									}
								}
							}
						} else {
							// частичный выкуп
							buyAltCount := buyQuantityAlt - orderQuantity
							worker.AddAlt(strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], buyAltCount, buyRate)
							worker.AvailableBTCCash -= buyAltCount * buyRate

							buyQuantityAlt = buyAltCount
						}
					}
				}

				if !findOrder {
					// выкупили полностью
					worker.AddAlt(strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], buyQuantityAlt, buyRate)
					worker.AvailableBTCCash -= buyQuantityAlt * buyRate

					uuidBuyOrder = ""
					worker.BuyActiveMarket = nil
					buyRate = 0
					buyQuantityAlt = 0
				}
			}
		}
	}
}