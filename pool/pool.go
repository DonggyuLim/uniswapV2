package pool

import (
	"log"
	"math/big"
)

type token struct {
	Name    string
	Symbol  string
	Balance *big.Int
}

func (t *token) GetTokenBalance() *big.Int {
	return t.Balance
}

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
	return x.Mul(x, y)
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

// business logic
func (p *Pool) Deposit(tokenA, tokenB token) token {

	x, y := p.Reserve() //pool reserve

	xDeposit, yDeposit := tokenA.Balance, tokenB.Balance
	if xDeposit.Cmp(yDeposit) == 1 {
		log.Panicf("%s 와 %s의 비율이 맞지 않습니다. %s가 더 작아야합니다", xDeposit, yDeposit, xDeposit)
	}
	//1:10
	//poolRatio = 10
	var poolRatio *big.Int
	poolRatio = poolRatio.Div(x, y)
	var depositRatio *big.Int
	//10:30
	//depositRatio = 3
	depositRatio = depositRatio.Div(xDeposit, yDeposit)
	pc := p.getPoolCoinBalance()

	var lp *big.Int
	//pr = 1:10  10
	//dr = 1:3    3
	switch poolRatio.Cmp(depositRatio) {
	//poolRatio > depositRatio
	case -1:
		return token{
			Name:    p.getPoolCoinName(),
			Balance: big.NewInt(0),
		}
	//poolRatio == depositRatio
	case 0:
		//lp = sqrt(x*y)
		//0이라면 풀의 비율과 보낸 토큰비율이 맞으니 그대로 저장하고 lp 토큰 발행
		lp = lp.Sqrt(xDeposit.Mul(xDeposit, yDeposit))
		p.Rs.X.Balance = x.Add(x, xDeposit)
		p.Rs.Y.Balance = y.Add(y, yDeposit)
		p.PC.Balance = pc.Add(pc, lp)
		return token{
			Name:    p.getPoolCoinName(),
			Balance: lp,
		}
	//poolRatio < depositRatio
	case 1:
		//lp = sqrt(x*y)
		var resultY *big.Int
		resultY = resultY.Mul(xDeposit, poolRatio)

		lp = lp.Sqrt(xDeposit.Mul(xDeposit, resultY))
		p.Rs.X.Balance = x.Add(x, xDeposit)
		p.Rs.Y.Balance = y.Add(x, y)
		p.PC.Balance = pc.Add(pc, lp)
		return token{
			Name:    p.getPoolCoinName(),
			Balance: lp,
		}
	default:
		return token{
			Name:    p.getPoolCoinName(),
			Balance: big.NewInt(0),
		}
	}
	//pool 비율과 보낸 토큰 비율이 다르면?pool 비율에 맞게 끔 토큰 비율만큼만 deposit
}

func (p *Pool) Withdraw(lp token) (x, y token) {
	if lp.Name != p.getPoolCoinName() {
		log.Panic("not pool coin")
	}
	poolBalance := p.getPoolCoinBalance()

	// (lp / poolBalance) *100
	// ex
	// pool.balance = 1000
	// lp = 10
	// percent = 1
	// (10/1000) * 100 = 1
	percent := lp.Balance.Mul(lp.Balance.Div(lp.Balance, poolBalance), big.NewInt(100))
	xBalance, yBalance := p.Reserve()
	hundread := big.NewInt(100)

	//xRefund = x*percent / 100
	//ex
	//x = 1000
	//percent = 1
	//(x * percent) / 100 = 10
	xRefund := hundread.Div(xBalance.Mul(xBalance, percent), hundread)

	yRefund := hundread.Div(yBalance.Mul(yBalance, percent), hundread)
	xName, yName := p.getPairNameFromPool()
	x = token{
		Name:    xName,
		Balance: xRefund,
	}
	y = token{
		Name:    yName,
		Balance: yRefund,
	}
	p.Rs.X.Balance = xBalance.Sub(xBalance, xRefund)
	p.Rs.Y.Balance = yBalance.Sub(yBalance, yRefund)
	p.PC.Balance = poolBalance.Sub(poolBalance, lp.Balance)
	return
}
