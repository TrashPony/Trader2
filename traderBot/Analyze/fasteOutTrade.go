package Analyze

import (
	"fmt"
	"github.com/TrashPony/Trader2/traderInfo"
	"github.com/TrashPony/Trader2/utils"
	"strconv"
)

const emergencyExitPercent = -7
const profitExit = 1

func BaseOutAlgorithm(market *traderInfo.Market, startProfit, GrowProfitPrice float64, buy bool) (Sell bool, fast bool, newProfit float64, Ask float64) {
	Ask, _ = market.MarketSummary.Ask.Float64()
	Bid, _ := market.MarketSummary.Bid.Float64()

	if Ask > GrowProfitPrice {
		startDifference := utils.PercentageCalculator(GrowProfitPrice, Ask)
		newProfit = Ask
		fmt.Println("СП: ", GrowProfitPrice, " НП: ", newProfit, " up ", startDifference, " % ")
		return false, false, newProfit, Ask
	}

	if Ask <= GrowProfitPrice {

		difference := utils.PercentageCalculator(GrowProfitPrice, Ask)
		startDifference := utils.PercentageCalculator(startProfit, Ask)

		if difference < emergencyExitPercent {
			fmt.Println("Цена упала на - "+strconv.Itoa(emergencyExitPercent)+"%, Экстренный перезакуп!!!", utils.PercentageCalculator(GrowProfitPrice, Ask))
			return true, true, GrowProfitPrice, Ask
		}

		if startDifference > 1 && !buy {
			fmt.Println("Алгоритм посчитал что рынок больше не эффективен")
			return true, false, GrowProfitPrice, Ask
		}
	}

	profitExitDifference := utils.PercentageCalculator(startProfit, Bid)
	if profitExitDifference > profitExit {
		fmt.Println("Алгоритм посчитал что прибыль достаточна и можно выходить")
		return true, true, GrowProfitPrice, Ask
	}

	return false, false, GrowProfitPrice, Ask
}
