package traderInfo

import (
	"github.com/toorop/go-bittrex"
	"strings"
)

func GetMarket(MarketName string) *Market {
	marketSummary, err := GetBittrex().GetMarketSummary(MarketName)
	if err != nil {
		println(err.Error())
		return nil
	}

	ordersBuy, err := GetBittrex().GetOrderBookBuySell(MarketName, "buy")
	if err != nil {
		println(err.Error())
		return nil
	}

	ordersSell, err := GetBittrex().GetOrderBookBuySell(MarketName, "sell")
	if err != nil {
		println(err.Error())
		return nil
	}

	marketHistory, err := GetBittrex().GetMarketHistory(MarketName)
	if err != nil {
		println(err.Error())
		return nil
	}

	myOpenOrders, err := GetBittrex().GetOpenOrders(MarketName)
	if err != nil {
		println(err.Error())
		return nil
	}

	market := Market{CurrencyPair: MarketName, MarketSummary: marketSummary[0], OrdersBuy: ordersBuy,
		OrdersSell: ordersSell, MarketHistory: marketHistory, MyOpenOrders: myOpenOrders}

	return &market
}

func GetAllMarket() []bittrex.Market {
	markets, err := GetBittrex().GetMarkets()
	if err != nil {
	}

	allAvailableMarket := make([]bittrex.Market, 0)

	for _, market := range markets {

		pair := strings.Split(market.MarketName, "-")

		if pair[0] == "BTC" && market.IsActive {
			if pair[1] == "BCH" ||
				pair[1] == "XMR" ||
				pair[1] == "ETH" ||
				pair[1] == "RDD" ||
				pair[1] == "OCN" {
				allAvailableMarket = append(allAvailableMarket, market)
			}
		}
	}

	return allAvailableMarket
}