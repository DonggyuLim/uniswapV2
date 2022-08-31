package pool

import (
	"github.com/shopspring/decimal"
)

// return x/y or y/x or 1

type reserve struct {
	X token
	Y token
}

type Pool struct {
	Rs reserve
	PC token
}

func CreatePool(tokenA, tokenB, pc token) *Pool {
	return &Pool{
		Rs: reserve{
			X: tokenA,
			Y: tokenB,
		},
		PC: pc,
	}
}

// return reserved coin
func (p *Pool) Reserve() (x, y decimal.Decimal) {
	x = p.Rs.X.Balance
	y = p.Rs.Y.Balance
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
func (p *Pool) getPoolCoinBalance() decimal.Decimal {
	return p.PC.Balance
}

// return pc name
func (p *Pool) getPoolCoinName() string {
	return p.PC.Name
}

// return reserve coin name
func (p *Pool) getPairNameFromPool() (string, string) {
	reserve := p.Rs
	return reserve.X.Name, reserve.Y.Name
}
