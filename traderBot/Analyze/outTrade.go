package Analyze

import (
	"../../traderInfo"
)

type AnalyzerOutTrade struct {
	Name string `json:"name"`
}

func (analyze *AnalyzerOutTrade) Analyze(market *traderInfo.Market, startProfit, GrowProfitPrice float64) (Sell bool, fast bool, newProfit float64, Ask float64) {

	if analyze.Name == "fast" {
		return BaseOutAlgorithm(market, startProfit, GrowProfitPrice)
	}

	return false, false, GrowProfitPrice, GrowProfitPrice
}
