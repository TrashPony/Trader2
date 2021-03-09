package Analyze

import (
	"github.com/TrashPony/Trader2/traderInfo"
)

type AnalyzerOutTrade struct {
	Name string `json:"name"`
}

func (analyze *AnalyzerOutTrade) Analyze(market *traderInfo.Market, startProfit, GrowProfitPrice float64, buy bool) (Sell bool, fast bool, newProfit float64, Ask float64) {

	if analyze.Name == "fast" {
		return BaseOutAlgorithm(market, startProfit, GrowProfitPrice, buy)
	}

	return false, false, GrowProfitPrice, GrowProfitPrice
}
