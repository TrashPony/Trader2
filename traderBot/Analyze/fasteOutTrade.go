package Analyze

import (
	"../../traderInfo"
	"../../utils"
	"fmt"
)

func BaseOutAlgorithm(market *traderInfo.Market, startProfit, GrowProfitPrice float64) (Sell bool, fast bool, newProfit float64, Ask float64) {
	Ask, _ = market.MarketSummary.Ask.Float64()

	if Ask > GrowProfitPrice {
		startDifference := utils.PercentageCalculator(GrowProfitPrice, Ask)
		newProfit = Ask
		fmt.Println("СП: ", GrowProfitPrice, " НП: ", newProfit, " up ", startDifference, " % ")
		return false, false, newProfit, Ask
	}

	if Ask <= GrowProfitPrice {
		difference := utils.PercentageCalculator(GrowProfitPrice, Ask)
		startDifference := utils.PercentageCalculator(startProfit, Ask)

		if difference < -2 {
			fmt.Println("Цена упала на - 2%, Экстренный перезакуп!!!", utils.PercentageCalculator(GrowProfitPrice, Ask))
			return true, true, GrowProfitPrice, Ask
		}

		if startDifference > 0.2 {
			fmt.Println("Алгоритм посчитал что рынок больше не эффективен")
			return true, false, GrowProfitPrice, Ask
		}
	}

	return false, false, GrowProfitPrice, Ask
}
