package Analyze

import (
	"../../traderInfo"
	"../../utils"
)

func BaseInAlgorithm(market *traderInfo.Market) (buy bool, fast bool) {
	Bid, _ := market.MarketSummary.Bid.Float64()
	Low, _ := market.MarketSummary.Low.Float64()
	High, _ := market.MarketSummary.High.Float64()
	FirstOrderBuy, _ := market.OrdersBuy[1].Rate.Float64()
	SecondOrderBuy, _ := market.OrdersBuy[0].Rate.Float64()

	avgLowHigh := (Low + High) / 2
	second := utils.PercentageCalculator(SecondOrderBuy, FirstOrderBuy)
	summaryAsc, summaryBid := get25HistorySumPrice(market)

	sumCap := summaryBid > summaryAsc
	lastPriceCheck := Bid >= avgLowHigh // ордер на покупку ниже средней по суткам
	secondCheck := second < 0.10

	return lastPriceCheck && secondCheck && sumCap, false
}

func get25HistorySumPrice(market *traderInfo.Market) (summaryAsc, summaryBid float64) {
	for _, order := range market.MarketHistory {
		quantity, _ := order.Quantity.Float64()
		price, _ := order.Price.Float64()

		if order.OrderType == "BUY" {
			summaryBid += quantity * price
		}
		if order.OrderType == "SELL" {
			summaryAsc += quantity * price
		}
	}

	return summaryAsc, summaryBid
}
