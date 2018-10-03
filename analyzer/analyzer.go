package analyzer

import (
	"fmt"
	"github.com/toorop/go-bittrex"
)

func Analyzer(marketSummary []bittrex.MarketSummary, ordersBuy []bittrex.Orderb, ordersSell []bittrex.Orderb) {

	Bid, _ := marketSummary[0].Bid.Float64()
	Ask, _ := marketSummary[0].Ask.Float64()
	Low, _ := marketSummary[0].Low.Float64()
	High, _ := marketSummary[0].High.Float64()
	FirstOrderBuy, _ := ordersBuy[1].Rate.Float64()
	SecondOrderBuy, _ := ordersBuy[0].Rate.Float64()

	avgLowHigh := (Low + High) / 2
	differenceAskBind := percentageCalculator(Bid, Ask)
	demand := percentageCalculator(float64(marketSummary[0].OpenSellOrders), float64(marketSummary[0].OpenBuyOrders))
	second := percentageCalculator(FirstOrderBuy, SecondOrderBuy)

	fmt.Println(demand, differenceAskBind, avgLowHigh, second)
}

func percentageCalculator(a, b float64) (result float64) {
	result = 100 - (a * 100 / b)
	return result
}
