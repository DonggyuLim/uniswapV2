package pool

import (
	"github.com/shopspring/decimal"
)

type token struct {
	Name    string
	Balance decimal.Decimal
}

func (t *token) GetTokenBalance() decimal.Decimal {
	return t.Balance
}
