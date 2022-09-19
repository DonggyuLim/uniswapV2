package pool

import "github.com/shopspring/decimal"

type Token struct {
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
	Account string          `json:"address"`
}

func (t *Token) GetBalance() decimal.Decimal {
	return t.Balance
}

func (t *Token) GetAddress() string {
	return t.Account
}

func (t *Token) GetTokenName() string {
	return t.Name
}
