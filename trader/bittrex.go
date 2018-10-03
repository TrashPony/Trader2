package trader

import (
	"../analyzer"
	"github.com/toorop/go-bittrex"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func GetMarket(MarketName string) {
	// Bittrex client
	bittrexMarket := bittrex.New(API_KEY, API_SECRET)

	marketSummary, err := bittrexMarket.GetMarketSummary(MarketName)
	if err != nil {
		panic(err)
	}

	ordersBuy, err := bittrexMarket.GetOrderBookBuySell(MarketName, "buy")
	if err != nil {
		panic(err)
	}

	ordersSell, err := bittrexMarket.GetOrderBookBuySell(MarketName, "sell")
	if err != nil {
		panic(err)
	}

	analyzer.Analyzer(marketSummary, ordersBuy, ordersSell)
}
