package Analyze

import "../../traderInfo"

type AnalyzerOutTrade struct {
	Name string `json:"name"`
}

func (analyze *AnalyzerOutTrade) Analyze(market *traderInfo.Market) (Sell bool, fast bool) {
	return true, false
}
