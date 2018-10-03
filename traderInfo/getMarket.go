package traderInfo

import "strings"

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

func GetAllMarket() []*Market {
	// TODO перенести это в торгового бота что бы не ждать пока получаться все рынки
	markets, err := GetBittrex().GetMarkets()
	if err != nil {

	}

	allMarket := make([]*Market, 0)

	for _, market := range markets {

		pair := strings.Split(market.MarketName, "-")

		if pair[0] == "BTC" && market.IsActive {
			if pair[1] == "BCH" ||
				pair[1] == "XMR" ||
				pair[1] == "ETH" ||
				pair[1] == "RDD" ||
				pair[1] == "OCN" { // из за того что рабоатет крайне медленно пока определил 5 валют для торга
				allMarket = append(allMarket, GetMarket(market.MarketName)) // крайне медленно работает
			}
		}
	}

	return allMarket
}
