package pool

import (
	m "cosmossdk.io/math"
)

type Token struct {
	Name    string `json:"name"`
	Balance m.Uint `json:"balance"`
	Account string `json:"address"`
}

func (t *Token) GetBalance() m.Uint {
	return t.Balance
}

func (t *Token) GetAccount() string {
	return t.Account
}

func (t *Token) GetTokenName() string {
	return t.Name
}
