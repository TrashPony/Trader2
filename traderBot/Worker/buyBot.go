package Worker

import (
	"../../traderInfo"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func (worker *Worker) TradeBuyBot() {

	markets := traderInfo.GetAllMarket()

	uuidBuyOrder := ""
	var buyRate float64
	var buyQuantityAlt float64
	var profitPrice float64

	for {

		if worker.AvailableBTCCash >= 0.0005 {

			var err error

			if uuidBuyOrder == "" && worker.BuyActiveMarket == nil {

				for _, marketName := range markets {

					market := traderInfo.GetMarket(marketName.MarketName)

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
							profitPrice = worker.AvailableBTCCash / quantity

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

								firstRate, okRate := first.Rate.Float64()
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
										profitPrice = 0
									}
								}
							}
						} else {
							// частичный выкуп

							buyAltCount := buyQuantityAlt - orderQuantity
							buyQuantityAlt -= buyAltCount

							worker.AddAlt(strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], buyAltCount, buyRate, profitPrice)
							worker.AvailableBTCCash -= buyAltCount * buyRate
							worker.SellActiveMarkets[worker.BuyActiveMarket.CurrencyPair] = worker.BuyActiveMarket

							fmt.Print("Купил ", buyAltCount, strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], " по ", buyRate)
						}
					}
				}

				if !findOrder {
					// выкупили полностью
					worker.AddAlt(strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], buyQuantityAlt, buyRate, profitPrice)
					worker.AvailableBTCCash -= buyQuantityAlt * buyRate
					worker.SellActiveMarkets[worker.BuyActiveMarket.CurrencyPair] = worker.BuyActiveMarket

					fmt.Print("Купил ", buyQuantityAlt, strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], " по ", buyRate)

					uuidBuyOrder = ""
					worker.BuyActiveMarket = nil
					buyRate = 0
					buyQuantityAlt = 0
					profitPrice = 0
				}
			}
		}
	}
}
