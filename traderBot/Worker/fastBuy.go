package Worker

import (
	"github.com/TrashPony/Trader2/traderInfo"
	"github.com/TrashPony/Trader2/utils"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func (worker *Worker) FastTradeBuy() {
	for {
		time.Sleep(1 * time.Second) // без этих слипов система виснет

		if worker.AvailableBTCCash >= 0.0005 {

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
						var logRow string

						if fast { // если фаст значит в рынок надо войти как можно быстрее )
							priceByu, _ = market.OrdersSell[0].Rate.Float64()
							logRow = "Fast buy - " + market.CurrencyPair + " по " + utils.FloatToString(priceByu)
						} else {
							priceByu, _ = market.OrdersBuy[0].Rate.Float64()
							priceByu += 0.00000001 // наращиваем 1 сатоши что бы стать самым первым ордером в стакане
							logRow = "Slow buy - " + market.CurrencyPair + " по " + utils.FloatToString(priceByu)
						}

						fee := worker.AvailableBTCCash * worker.Fee
						quantity := (worker.AvailableBTCCash - fee) / priceByu

						minSize := traderInfo.Markets.GetMinTradeSize(market.MarketSummary.MarketName)
						if minSize > quantity {
							continue
						}

						worker.AddLog(logRow)

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

				priceRemaining, _ := worker.BuyOrder.Price.Float64() // это те деньни которые были потрачены на частичный выкуп

				profitPrice := utils.Round(GetFee(buyQuantityAlt, buyRate, worker.Fee), 8)

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

						buyAltCount := buyQuantityAlt - buyQuantityRemaining
						worker.AddAlt(strings.Split(market.CurrencyPair, "-")[1], buyAltCount, buyRate, profitPrice)

						// todo тут завышеная коммисия т.к. надо высчитывать по доли
						worker.AvailableBTCCash = worker.AvailableBTCCash - ((priceRemaining - priceRemaining) + fee)

						worker.AddLog("Купил частичку" + utils.FloatToString(buyAltCount) + " " + strings.Split(market.CurrencyPair, "-")[1] +
							" по " + utils.FloatToString(buyRate))
					}
				} else {

					worker.AddAlt(strings.Split(market.CurrencyPair, "-")[1], buyQuantityAlt, buyRate, profitPrice)
					//                                         (цена за коины) + (коммисия)
					worker.AvailableBTCCash = worker.AvailableBTCCash - (priceRemaining + fee)
					worker.BuyOrder = nil

					worker.AddLog("Купил " + utils.FloatToString(buyQuantityAlt) + " " + strings.Split(market.CurrencyPair, "-")[1] +
						" по " + utils.FloatToString(buyRate))
				}
			}
		}
	}
}

func GetFee(quantity, price, baseFee float64) float64 {

	/*
		baseFee = 0.0075
		q = 0.10560508
		p = 0.00575602

		f = q * p * baseFee = 0.00000455
		needPrice = q * p + f = 0.00061241
		needProfitPrice =  needPrice + (needPrice * (baseFee/100)) = 0.000617003075
		profitOut = needPrice / q = 0.00584268
	*/

	fee := quantity * price * baseFee
	needPrice := quantity*price + fee
	needProfitPrice := needPrice + (needPrice * baseFee)
	profitOut := needProfitPrice / quantity
	return profitOut
}
