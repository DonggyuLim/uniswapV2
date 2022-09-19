package pool

import (
	"github.com/DonggyuLim/uniswap/utils"
	"github.com/shopspring/decimal"
)

const feeRate = "0.0003"

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

type Pool struct {
	Name    string
	FeeRate decimal.Decimal
	XFee    decimal.Decimal
	YFee    decimal.Decimal
	X       Token
	Y       Token
	LP      poolToken
}

func CreatePool(tokenA, tokenB Token, lp poolToken) Pool {
	feeRate, _ := decimal.NewFromString(feeRate)
	zero := decimal.NewFromInt(0)
	return Pool{
		Name:    utils.GetKey(tokenA.GetTokenName(), tokenB.GetTokenName()),
		X:       tokenA,
		Y:       tokenB,
		XFee:    zero,
		YFee:    zero,
		FeeRate: feeRate,
		LP:      lp,
	}
}

func (p *Pool) GetName() string {
	return p.Name
}

func (p *Pool) GetXName() string {
	return p.X.Name
}
func (p *Pool) GetYName() string {
	return p.Y.Name
}
func (p *Pool) GetFeeRate() decimal.Decimal {
	return p.FeeRate
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
func (p *Pool) GetLPname() string {
	return p.LP.GetName()
}

// return reserve coin name
func (p *Pool) getPairNameFromPool() (string, string) {
	return p.X.Name, p.Y.Name
}
