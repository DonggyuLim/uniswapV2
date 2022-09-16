package pool

import (
	"github.com/shopspring/decimal"
)

const fee = "0.03"

// return x/y or y/x or 1
type Token struct {
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
	Address string          `json:"address"`
}

func (t *Token) GetBalance() decimal.Decimal {
	return t.Balance
}

func (t *Token) TokenName() string {
	return t.Name
}

type Pool struct {
	Name string
	Fee  decimal.Decimal
	X    Token
	Y    Token
	LP   poolToken
}

func CreatePool(tokenA, tokenB Token, lp poolToken) Pool {
	fee, _ := decimal.NewFromString(fee)
	return Pool{

		X:   tokenA,
		Y:   tokenB,
		Fee: fee,
		LP:  lp,
	}
}

func (p *Pool) GetName() string {
	return p.Name
}

// return reserved coin
func (p *Pool) Reserve() (x, y decimal.Decimal) {
	x = p.X.Balance
	y = p.Y.Balance
	return
}

// return x*y = k
func (p *Pool) K() decimal.Decimal {

	x, y := p.Reserve()
	z := x.Mul(y)
	return z
}

// return x/y = price
func (p *Pool) poolPrice() (price decimal.Decimal) {
	x, y := p.Reserve()
	price = x.Div(y)
	return
}

// return PC.balance
// func (p *Pool) getPoolCoinBalance() decimal.Decimal {
// 	return p.PC.Balance
// }

// return pc name
func (p *Pool) getLPname() string {
	return p.LP.GetName()
}

// return reserve coin name
func (p *Pool) getPairNameFromPool() (string, string) {
	return p.X.Name, p.Y.Name
}
