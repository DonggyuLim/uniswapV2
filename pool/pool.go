package pool

import (
	"log"
	"math/big"
)

type token struct {
	name    string
	balance *big.Int
}

func (t *token) GetTokenBalance() *big.Int {
	return t.balance
}

type poolCoin struct {
	name    string
	balance *big.Int
}

type reserve struct {
	X token
	Y token
}

type Pool struct {
	Rs reserve
	PC poolCoin
}

func CreatePool(tokenA, tokenB token, pc poolCoin) *Pool {
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
	x = p.Rs.X.balance
	y = p.Rs.Y.balance
	return
}

func (p *Pool) K() *big.Int {
	x, y := p.Reserve()
	return x.Mul(x, y)
}

// return poolCoin.balance
func (p *Pool) getPoolCoinBalance() *big.Int {
	return p.PC.balance
}

func (p *Pool) getPoolCoinName() string {
	return p.PC.name
}

// return reserve coin name
func (p *Pool) pairNameFromPool() (string, string) {
	reserve := p.Rs
	return reserve.X.name, reserve.Y.name
}

// business logic
func (p *Pool) Deposit(tokenA, tokenB token) poolCoin {
	x, y := p.Reserve() //pool reserve
	pc := p.getPoolCoinBalance()
	aBalance, bBlance := tokenA.balance, tokenB.balance
	// xInt, yInt := x.Int64(), y.Int64()
	poolRatio := x.Div(x, y)

	tokenRatio := aBalance.Div(aBalance, bBlance)
	//pool 비율과 보낸 토큰 비율이 다르면?pool 비율에 맞게 끔 토큰 비율만큼만 deposit
	//1:10 비율이라면?
	// 패닉을 하는게 아니라
	// 비율만큼 업데이트하고 나머지는 환불해줘야함.
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	if poolRatio != tokenRatio {
		log.Panic("To deposit is equal ratio which Rx div Ry  ")
	}
	pCoin := pc.Div(pc, poolRatio)
	p.Rs.X.balance = x.Add(x, aBalance)
	p.Rs.Y.balance = y.Add(y, bBlance)
	p.PC.balance = pc.Add(pc, pCoin)
	return poolCoin{
		name:    p.getPoolCoinName(),
		balance: pCoin,
	}
}
