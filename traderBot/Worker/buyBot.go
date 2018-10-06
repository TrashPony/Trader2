package Worker

import (
	"../../traderInfo"
	"../../utils"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func (worker *Worker) TradeBuyBot() {

	uuidBuyOrder := ""
	var buyRate float64
	var buyQuantityAlt float64
	var profitPrice float64

	for {
		if worker.AvailableBTCCash >= 0 { //.0005 {

			markets := traderInfo.Markets.GetMarketPool()

			var err error

			if uuidBuyOrder == "" && worker.BuyActiveMarket == nil {

				for _, market := range markets {

					err = market.UpdateMarket()
					if err != nil {
						continue
					}

					buy, fast := worker.InTradeStrategy.Analyze(market)

					if buy {

						var priceByu float64

						if fast { // если фаст значит в рынок надо войти как можно быстрее )
							priceByu, _ = market.OrdersSell[0].Rate.Float64()
							worker.AddLog("Fast buy - " + market.CurrencyPair)
						} else {
							priceByu, _ = market.OrdersBuy[0].Rate.Float64()
							priceByu += 0.00000001 // наращиваем 1 сатоши что бы стать самым первым ордером в стакане
							worker.AddLog("Slow buy - " + market.CurrencyPair)
						}

						fee := worker.AvailableBTCCash * worker.Fee
						quantity := (worker.AvailableBTCCash - fee) / priceByu

						uuidBuyOrder, err = market.BuyLimit(decimal.NewFromFloat(quantity), decimal.NewFromFloat(priceByu))
						if err != nil {
							worker.AddLog("Error buy - " + market.CurrencyPair + err.Error())
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

								firstRate, _ := first.Rate.Float64()
								buy, _ := worker.InTradeStrategy.Analyze(worker.BuyActiveMarket)

								if firstRate > buyRate && buy {
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

									worker.AddLog("Cancel buy - " + worker.BuyActiveMarket.CurrencyPair)
								}
							}
						} else {
							// частичный выкуп

							buyAltCount := buyQuantityAlt - orderQuantity
							buyQuantityAlt -= buyAltCount

							worker.AddAlt(strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], buyAltCount, buyRate, profitPrice)
							worker.AvailableBTCCash -= buyAltCount * buyRate

							worker.AddLog("Купил " + utils.FloatToString(buyAltCount) + " " + strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1] +
								" по " + utils.FloatToString(buyRate))
						}
					}
				}

				if !findOrder {
					// выкупили полностью
					worker.AddAlt(strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1], buyQuantityAlt, buyRate, profitPrice)
					worker.AvailableBTCCash -= buyQuantityAlt * buyRate

					worker.AddLog("Купил " + utils.FloatToString(buyQuantityAlt) + " " + strings.Split(worker.BuyActiveMarket.CurrencyPair, "-")[1] +
						" по " + utils.FloatToString(buyRate))

					uuidBuyOrder = ""
					worker.BuyActiveMarket = nil
					buyRate = 0
					buyQuantityAlt = 0
					profitPrice = 0
				}
			}
		}
		time.Sleep(1 * time.Second) // без этих слипов система виснет
	}
}
