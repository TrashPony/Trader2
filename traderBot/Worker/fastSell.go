package Worker

import (
	"../../traderInfo"
	"../../utils"
	"github.com/shopspring/decimal"
	"time"
)

func (worker *Worker) FastTradeSell() {
	for {
		for _, altBalances := range worker.AltBalances { // берем прошлые сделки покупок
			if !altBalances.Sell && altBalances.Balance*altBalances.ProfitPrice > 0.0006 {
				sellAltCoin(worker, altBalances)
			}
		}
		time.Sleep(1 * time.Second) // без этих слипов система виснет
	}
}

func sellAltCoin(worker *Worker, altBalances *Alt) {

	market := traderInfo.GetMarket("BTC-" + altBalances.AltName)
	if market == nil {
		return
	}

	if altBalances.SellOrder == nil {
		sell, fast, newProfit, asc := worker.OutTradeStrategy.Analyze(market, altBalances.ProfitPrice, altBalances.GrowProfitPrice)

		altBalances.GrowProfitPrice = newProfit
		altBalances.TopAsc = asc

		if sell {
			var priceSell float64

			if fast { // если фаст значит из рынка нужно выйти немедленно
				priceSell, _ = market.OrdersBuy[0].Rate.Float64()
				worker.AddLog("Fast sell - " + market.CurrencyPair + " за " + utils.FloatToString(priceSell))
			} else {
				priceSell, _ = market.OrdersSell[0].Rate.Float64()
				priceSell -= 0.00000001 // уменьшаем на 1 сатоши что бы стать самым первым ордером в стакане
				worker.AddLog("Slow sell - " + market.CurrencyPair + " за " + utils.FloatToString(priceSell))
			}

			uuidSellOrder, err := market.SellLimit(decimal.NewFromFloat(altBalances.Balance), decimal.NewFromFloat(priceSell))
			if err != nil {
				println(err.Error())
				worker.AddLog("Error sell - " + market.CurrencyPair + " " + err.Error())
			}

			altBalances.SellOrder, _ = market.GetOrder(uuidSellOrder)
		}
	} else {
		sellOrder, err := market.GetOrder(altBalances.SellOrder.OrderUuid)
		if err != nil || sellOrder == nil {
			return
		}
		altBalances.SellOrder = sellOrder

		sellQuantityAlt, _ := altBalances.SellOrder.Quantity.Float64()
		sellQuantityRemaining, _ := altBalances.SellOrder.QuantityRemaining.Float64()
		sellRate, _ := altBalances.SellOrder.Limit.Float64()
		fee, _ := altBalances.SellOrder.CommissionReserved.Float64() // это комисиия в оредере

		if sellOrder.IsOpen {
			if sellQuantityAlt == sellQuantityRemaining {
				// не выкупили
				first, err := market.GetFirstSellOrder()
				if err != nil {
					return
				} else {

					firstRate, _ := first.Rate.Float64()
					sell, _, newProfit, asc := worker.OutTradeStrategy.Analyze(market, altBalances.ProfitPrice, altBalances.GrowProfitPrice)
					altBalances.GrowProfitPrice = newProfit
					altBalances.TopAsc = asc

					if firstRate < sellRate || !sell {
						// если у первого заказа цена выше чем у наc
						// и рынок досихпор считается перспективным то пересоздаем ордер
						err = market.CancelOrder(altBalances.SellOrder.OrderUuid)
						if err != nil {
							println(err.Error())
							return
						} else {
							altBalances.SellOrder = nil
							worker.AddLog("Cancel sell - " + market.CurrencyPair)
						}
					}
				}
			} else {
				// частичный выкуп
				sellAltCount := sellQuantityAlt - sellQuantityRemaining    // высчитываем сколько купили
				worker.AvailableBTCCash += (sellAltCount * sellRate) - fee // высчитывает заработаные BTC
				altBalances.Balance -= sellAltCount                        // отнимать у баланса валюты проданые монеты

				worker.AddLog("Продалась часть закупа с выгодой " +
					utils.FloatToString(utils.PercentageCalculator(altBalances.ProfitPrice, sellRate)))
			}
		} else {
			// выкупили полностью
			worker.AvailableBTCCash += (altBalances.Balance * sellRate) - fee
			sellOrder = nil

			altBalances.Sell = true
			altBalances.SellRate = sellRate

			worker.AddLog("Продалось относительно начального закупа с выгодой " +
				utils.FloatToString(utils.PercentageCalculator(altBalances.ProfitPrice, sellRate)))
		}
	}
}
