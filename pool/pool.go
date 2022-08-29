package pool

import (
	"math/big"
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
func (p *Pool) Reserve() (x, y *big.Int) {
	x = p.Rs.X.Balance
	y = p.Rs.Y.Balance
	return
}

// return x*y = k
func (p *Pool) K() *big.Int {

	x, y := p.Reserve()
	k := big.Int{}
	return k.Mul(x, y)
}

// return PC.balance
func (p *Pool) getPoolCoinBalance() *big.Int {
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
