package Analyze

import (
	"../../traderInfo"
)

type AnalyzerInTrade struct {
	Name string `json:"name"`
}

func (analyze *AnalyzerInTrade) Analyze(market *traderInfo.Market) (buy bool, fast bool) {

	if analyze.Name == "fast" {
		return BaseInAlgorithm(market)
	}

	return false, false
}
