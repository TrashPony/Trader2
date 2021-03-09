package Analyze

import (
	"github.com/TrashPony/Trader2/traderInfo"
	"github.com/TrashPony/Trader2/utils"
	"github.com/toorop/go-bittrex"
	"sort"
)

func BaseInAlgorithm(market *traderInfo.Market) (buy bool, fast bool) {
	Bid, _ := market.MarketSummary.Bid.Float64()
	//Low, _ := market.MarketSummary.Low.Float64()
	//High, _ := market.MarketSummary.High.Float64()
	FirstOrderBuy, _ := market.OrdersBuy[1].Rate.Float64()
	SecondOrderBuy, _ := market.OrdersBuy[0].Rate.Float64()

	//avgLowHigh := (Low + High) / 2
	second := utils.PercentageCalculator(SecondOrderBuy, FirstOrderBuy)
	summaryHistoryAsc, summaryHistoryBid := get25HistorySumPrice(market)

	summaryAsc, summaryBid := get25SumPrice(market)

	sumHistoryCap := summaryHistoryBid > summaryHistoryAsc // в истории преобладает покупки валюты по объему
	sumCap := summaryAsc > summaryBid                      // текущие ордера преобладают по объему на покупку валюты
	//lastPriceCheck := Bid >= avgLowHigh 				   // ордер на покупку ниже средней по суткам
	secondCheck := second < 0.10 // защита от тролинга
	bigPrice := Bid > 0.00005

	return secondCheck && sumHistoryCap && sumCap && bigPrice, false
}

func get25HistorySumPrice(market *traderInfo.Market) (summaryAsc, summaryBid float64) {

	countOrder := 0

	for _, order := range sortHistoryOrders(market) {

		quantity, _ := order.Quantity.Float64()
		price, _ := order.Price.Float64()
		countOrder++

		if countOrder < 25 && order.OrderType == "SELL" {
			summaryAsc += quantity * price
		}

		if countOrder < 25 && order.OrderType == "BUY" {
			summaryBid += quantity * price
		}
	}

	return summaryAsc, summaryBid
}

func get25SumPrice(market *traderInfo.Market) (summaryAsc, summaryBid float64) {
	countSell := 0
	countBuy := 0

	for _, order := range sortOrders("SELL", market.OrdersSell) {

		quantity, _ := order.Quantity.Float64()
		price, _ := order.Rate.Float64()

		if countSell < 25 {
			countSell++
			summaryAsc += quantity * price
		}
	}

	for _, order := range sortOrders("BUY", market.OrdersBuy) {

		quantity, _ := order.Quantity.Float64()
		price, _ := order.Rate.Float64()

		if countBuy < 25 {
			countBuy++
			summaryBid += quantity * price
		}
	}

	return summaryAsc, summaryBid
}

func sortOrders(typeOrders string, marketOrders []bittrex.Orderb) []bittrex.Orderb {

	sort.SliceStable(marketOrders, func(i, j int) bool {
		price1, _ := marketOrders[i].Rate.Float64()
		price2, _ := marketOrders[j].Rate.Float64()

		if typeOrders == "SELL" {
			return price1 < price2
		} else {
			return price1 > price2
		}
	})

	return marketOrders
}

func sortHistoryOrders(market *traderInfo.Market) []bittrex.Trade {
	ordersArray := make([]bittrex.Trade, 0)

	for _, order := range market.MarketHistory {
		ordersArray = append(ordersArray, order)
	}

	sort.SliceStable(ordersArray, func(i, j int) bool {
		time1 := ordersArray[i].Timestamp.Unix()
		time2 := ordersArray[j].Timestamp.Unix()

		return time1 > time2
	})

	return ordersArray
}
