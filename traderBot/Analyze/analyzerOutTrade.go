package Analyze

import (
	"../../traderInfo"
	"../../utils"
	"fmt"
)

type AnalyzerOutTrade struct {
	Name string `json:"name"`
}

func (analyze *AnalyzerOutTrade) Analyze(market *traderInfo.Market, startProfit, GrowProfitPrice float64) (Sell bool, fast bool, newProfit float64, Ask float64) {
	if analyze.Name == "BaseOutAlgorithm" {
		return BaseOutAlgorithm(market, startProfit, GrowProfitPrice)
	}

	return false, false, GrowProfitPrice, GrowProfitPrice
}

func BaseOutAlgorithm(market *traderInfo.Market, startProfit, GrowProfitPrice float64) (Sell bool, fast bool, newProfit float64, Ask float64) {
	Ask, _ = market.MarketSummary.Ask.Float64()

	if Ask > GrowProfitPrice {
		startDifference := utils.PercentageCalculator(GrowProfitPrice, Ask)
		newProfit = Ask
		fmt.Print("СП: ", GrowProfitPrice, " НП: ", newProfit, " up ", startDifference, " % ")
		return false, false, newProfit, Ask
	}

	if Ask <= GrowProfitPrice {
		difference := utils.PercentageCalculator(GrowProfitPrice, Ask)
		startDifference := utils.PercentageCalculator(startProfit, Ask)

		if difference < -2 {
			fmt.Print("Цена упала на - 2%, Экстренный перезакуп!!!", utils.PercentageCalculator(GrowProfitPrice, Ask))
			return true, true, GrowProfitPrice, Ask
		}

		if startDifference > 0.3 {
			fmt.Print("Алгоритм посчитал что рынок больше не эффективен")
			return true, false, GrowProfitPrice, Ask
		}
	}

	return false, false, GrowProfitPrice, Ask
}
