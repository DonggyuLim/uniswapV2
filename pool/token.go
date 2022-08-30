package pool

import "math/big"

type token struct {
	Name    string
	Balance *big.Float
}

func (t *token) GetTokenBalance() *big.Float {
	return t.Balance
}
