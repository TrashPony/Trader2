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

func (market *Market) SellLimit(quantity, rate decimal.Decimal) (string, error) {
	uuid, err := GetBittrex().SellLimit(market.CurrencyPair, quantity, rate)
	return uuid, err
}

func (market *Market) GetOpenOrders() ([]bittrex.Order, error) {
	orders, err := GetBittrex().GetOpenOrders(market.CurrencyPair)
	return orders, err
}

func (market *Market) GetOrder(uuid string) (*bittrex.Order2, error) {
	order, err := GetBittrex().GetOrder(uuid)
	return &order, err
}

func (market *Market) CancelOrder(uuid string) error {
	err := GetBittrex().CancelOrder(uuid)
	return err
}

func (market *Market) GetUserOrderHistory() ([]bittrex.Order, error) {
	orderHistory, err := GetBittrex().GetOrderHistory(market.CurrencyPair)
	return orderHistory, err
}

func (market *Market) GetFirstBuyOrder() (*bittrex.Orderb, error) {
	ordersBuy, err := GetBittrex().GetOrderBookBuySell(market.CurrencyPair, "buy")
	if len(ordersBuy) == 0 {
		println(err.Error())
		return nil, err
	}
	return &ordersBuy[0], err
}

func (market *Market) GetFirstSellOrder() (bittrex.Orderb, error) {
	ordersSell, err := GetBittrex().GetOrderBookBuySell(market.CurrencyPair, "sell")
	return ordersSell[0], err
}
