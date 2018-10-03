package traderInfo

import (
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"
)

type Account struct {
	Balances     []bittrex.Balance `json:"balances"`
	OrderHistory []bittrex.Order   `json:"order_history"`
}

func GetAccount() *Account {
	balances, err := GetBittrex().GetBalances()
	if err != nil {
		return nil
	}

	orderHistory, err := GetBittrex().GetOrderHistory("BTC-DOGE")
	if err != nil {
		return nil
	}

	account := Account{Balances: balances, OrderHistory: orderHistory}

	return &account
}

func (account *Account) GetAvailableCurrencyBalance(currencyName string) *decimal.Decimal {
	balance, err := GetBittrex().GetBalance(currencyName)
	if err != nil {
		return &balance.Available
	}

	return nil
}

func (account *Account) UpdateAccount() error {
	balances, err := GetBittrex().GetBalances()
	if err != nil {
		return err
	}

	orderHistory, err := GetBittrex().GetOrderHistory("BTC-DOGE")
	if err != nil {
		return err
	}

	account.Balances = balances
	account.OrderHistory = orderHistory

	return nil
}
