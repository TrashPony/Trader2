package Worker

import (
	"../../traderInfo"
	"../../utils"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func (worker *Worker) FastTradeBuy() {
	for {
		if worker.AvailableBTCCash >= 0.0006 {

			if worker.BuyOrder == nil {

				markets := traderInfo.Markets.GetMarketPool()

				for _, market := range markets {

					employment := CheckEmploymentMarket(market.CurrencyPair)
					if employment {
						continue // рынок уже кто то обрабатывает
					}

					err := market.UpdateMarket()
					if err != nil {
						continue
					}

					worker.BuyActiveMarket = market.CurrencyPair
					buy, fast := worker.InTradeStrategy.Analyze(market)

					if buy {

						var priceByu float64

						if fast { // если фаст значит в рынок надо войти как можно быстрее )
							priceByu, _ = market.OrdersSell[0].Rate.Float64()
							worker.AddLog("Fast buy - " + market.CurrencyPair + " по " + utils.FloatToString(priceByu))
						} else {
							priceByu, _ = market.OrdersBuy[0].Rate.Float64()
							priceByu += 0.00000001 // наращиваем 1 сатоши что бы стать самым первым ордером в стакане
							worker.AddLog("Slow buy - " + market.CurrencyPair + " по " + utils.FloatToString(priceByu))
						}

						fee := worker.AvailableBTCCash * worker.Fee
						quantity := (worker.AvailableBTCCash - fee) / priceByu

						uuidBuyOrder, err := market.BuyLimit(decimal.NewFromFloat(quantity), decimal.NewFromFloat(priceByu))
						if err != nil {
							worker.AddLog("Error buy - " + market.CurrencyPair + " " + err.Error())
							println(err.Error())
							continue
						}

						worker.BuyOrder, _ = market.GetOrder(uuidBuyOrder)
						break
					}
				}
			} else {

				market := traderInfo.GetMarket(worker.BuyActiveMarket)
				if market == nil {
					continue
				}

				buyOrder, err := market.GetOrder(worker.BuyOrder.OrderUuid)
				if err != nil || buyOrder == nil {
					continue
				}
				worker.BuyOrder = buyOrder

				buyQuantityAlt, _ := worker.BuyOrder.Quantity.Float64()                // количество монет
				buyQuantityRemaining, _ := worker.BuyOrder.QuantityRemaining.Float64() // количество оставшихся некупленных монет
				buyRate, _ := worker.BuyOrder.Limit.Float64()                          // цена покупки
				fee, _ := worker.BuyOrder.CommissionReserved.Float64()                 // это комисиия в оредере

				price, _ := worker.BuyOrder.Reserved.Float64()       // это сумарная стоимость всех монет в ордере без комисии
				priceRemaining, _ := worker.BuyOrder.Price.Float64() // это те деньни которые были потрачены на частичный выкуп

				if worker.BuyOrder.IsOpen {
					if buyQuantityAlt == buyQuantityRemaining {
						// не выкупили

						first, err := market.GetFirstBuyOrder()
						if err != nil {
							continue
						} else {

							firstRate, _ := first.Rate.Float64()
							buy, _ := worker.InTradeStrategy.Analyze(market)

							if firstRate > buyRate || !buy {
								// если у первого заказа цена выше чем у наc или рынок не перспективен то закрываем ордер
								err = market.CancelOrder(worker.BuyOrder.OrderUuid)
								if err != nil {
									println(err.Error())
								} else {
									worker.BuyOrder = nil
								}

								worker.AddLog("Cancel buy - " + market.CurrencyPair)
							}
						}
					} else {
						// частичный выкуп
						//                 ((цена за все коины) + (коммисия)) / (количество монет)
						profitPrice := utils.Round((worker.AvailableBTCCash+fee)/buyQuantityAlt, 8)

						buyAltCount := buyQuantityAlt - buyQuantityRemaining
						worker.AddAlt(strings.Split(market.CurrencyPair, "-")[1], buyAltCount, buyRate, profitPrice)

						// todo тут завышеная коммисия т.к. надо высчитывать по доли
						worker.AvailableBTCCash = worker.AvailableBTCCash - ((price - priceRemaining) + fee)

						worker.AddLog("Купил частичку" + utils.FloatToString(buyAltCount) + " " + strings.Split(market.CurrencyPair, "-")[1] +
							" по " + utils.FloatToString(buyRate))
					}
				} else {
					// выкупили полностью
					//                 ((цена за все коины) + (коммисия)) / (количество монет)
					profitPrice := utils.Round((worker.AvailableBTCCash+fee)/buyQuantityAlt, 8)

					worker.AddAlt(strings.Split(market.CurrencyPair, "-")[1], buyQuantityAlt, buyRate, profitPrice)
					//                                         (цена за коины) + (коммисия)
					worker.AvailableBTCCash = worker.AvailableBTCCash - (price + fee)
					worker.BuyOrder = nil

					worker.AddLog("Купил " + utils.FloatToString(buyQuantityAlt) + " " + strings.Split(market.CurrencyPair, "-")[1] +
						" по " + utils.FloatToString(buyRate))
				}
			}
		}
		time.Sleep(1 * time.Second) // без этих слипов система виснет
	}
}
