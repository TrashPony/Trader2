package traderInfo

import (
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"
)

type Market struct {
	CurrencyPair  string                `json:"currency_pair"`
	MarketSummary bittrex.MarketSummary `json:"market_summary"`
	OrdersBuy     []bittrex.Orderb      `json:"orders_buy"`
	OrdersSell    []bittrex.Orderb      `json:"orders_sell"`
	MarketHistory []bittrex.Trade       `json:"market_history"`
	MyOpenOrders  []bittrex.Order       `json:"my_open_orders"`
}

func (market *Market) UpdateMarket() error {
	marketSummary, err := GetBittrex().GetMarketSummary(market.CurrencyPair)
	if err != nil {
		return nil
	}

	ordersBuy, err := GetBittrex().GetOrderBookBuySell(market.CurrencyPair, "buy")
	if err != nil {
		return nil
	}

	ordersSell, err := GetBittrex().GetOrderBookBuySell(market.CurrencyPair, "sell")
	if err != nil {
		return nil
	}

	marketHistory, err := GetBittrex().GetMarketHistory(market.CurrencyPair)
	if err != nil {
		return nil
	}

	myOpenOrders, err := GetBittrex().GetOpenOrders(market.CurrencyPair)
	if err != nil {
		return nil
	}

	market.MarketSummary = marketSummary[0]
	market.OrdersBuy = ordersBuy
	market.OrdersSell = ordersSell
	market.MarketHistory = marketHistory
	market.MyOpenOrders = myOpenOrders

	return nil
}

func (market *Market) BuyLimit(quantity, rate decimal.Decimal) (string, error) {
	uuid, err := GetBittrex().BuyLimit(market.CurrencyPair, quantity, rate)
	return uuid, err
}

func (market *Market) SellLimit(quantity, rate *decimal.Decimal) (string, error) {
	uuid, err := GetBittrex().SellLimit(market.CurrencyPair, *quantity, *rate)
	return uuid, err
}

func (market *Market) GetOpenOrders() ([]bittrex.Order, error) {
	orders, err := GetBittrex().GetOpenOrders(market.CurrencyPair)
	return orders, err
}

func (market *Market) CancelOrder(uuid string) error {
	err := GetBittrex().CancelOrder(uuid)
	return err
}
