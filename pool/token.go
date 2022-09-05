package pool

import (
	"github.com/shopspring/decimal"
)

type token struct {
	Name    string
	Balance decimal.Decimal
}

func (t *token) TokenBalance() decimal.Decimal {
	return t.Balance
}

func (t *token) TokenName() string {
	return t.Name
}
