package traderInfo

import (
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"
)

type Account struct {
	Balances []bittrex.Balance `json:"balances"`
	StartBTC float64           `json:"start_btc"`
}

func GetAccount() *Account {
	balances, err := GetBittrex().GetBalances()
	if err != nil {
		println(err.Error())
		return nil
	}

	account := Account{Balances: balances}

	return &account
}

func (account *Account) GetAvailableCurrencyBalance(currencyName string) *decimal.Decimal {
	balance, err := GetBittrex().GetBalance(currencyName)
	if err != nil {
		println(err.Error())
		return nil
	}

	return &balance.Available
}

func (account *Account) UpdateAccount() error {
	balances, err := GetBittrex().GetBalances()
	if err != nil {
		println(err.Error())
		return err
	}

	account.Balances = balances

	return nil
}
