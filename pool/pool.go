package pool

import (
	"errors"

	m "cosmossdk.io/math"
	"github.com/DonggyuLim/uniswap/utils"
)

const feeRate = "0.3"

type Pool struct {
	Name    string
	FeeRate m.LegacyDec
	XFee    m.Uint
	YFee    m.Uint
	X       Token
	Y       Token
	LP      poolToken
}

func CreatePool(tokenA, tokenB Token, lp poolToken) Pool {
	feeRate := m.LegacyMustNewDecFromStr(feeRate)

	return Pool{
		Name:    utils.GetKey(tokenA.GetTokenName(), tokenB.GetTokenName()),
		X:       tokenA,
		Y:       tokenB,
		XFee:    m.NewUint(0),
		YFee:    m.NewUint(0),
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
func (p *Pool) GetFeeRate() m.LegacyDec {
	return p.FeeRate
}

// return reserved coin
func (p *Pool) Reserve() (x, y m.Uint) {
	x = p.X.Balance
	y = p.Y.Balance
	return
}

// return x*y = k
func (p *Pool) K() m.Uint {

	x, y := p.Reserve()

	z := x.Mul(y)
	return z
}

// return x/y = price
func (p *Pool) poolPrice() (price m.LegacyDec) {
	x, y := p.Reserve()
	dx, dy := m.LegacyMustNewDecFromStr(x.String()), m.LegacyMustNewDecFromStr(y.String())
	price = dx.Quo(dy)
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

// balance of account compare amount
func (p *Pool) lpCheckBalance(account string, amount m.Uint) error {
	balance := p.LP.BalanceOf(account)
	if balance.LT(amount) {
		return errors.New("you have not enough lp")
	}
	return nil
}
