package pool

import "math/big"

type token struct {
	Name    string
	Balance *big.Int
}

func (t *token) GetTokenBalance() *big.Int {
	return t.Balance
}
