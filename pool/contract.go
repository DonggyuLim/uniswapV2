package pool

import (
	"fmt"
	"log"
	"math/big"
)

// business logic
// 난중에DB 업데이트해야함.
func (p *Pool) Deposit(tokenA, tokenB token) (lp *token) {

	x, y := p.Reserve() //pool reserve

	xDeposit, yDeposit := tokenA.Balance, tokenB.Balance
	fmt.Printf("xDeposit = %v , yDeposit =%v\n", xDeposit, yDeposit)
	poolPrice, err := Price(x, y)

	if err != nil {
		log.Print(err)
	}

	fmt.Printf("poolPrice = %v\n", poolPrice)

	depositPrice, err := Price(xDeposit, yDeposit)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("depositPrice = %v\n", depositPrice)

	pc := p.getPoolCoinBalance()

	switch poolPrice.Cmp(depositPrice) {
	//poolPrice < depositPrice
	case -1:
		// x = 30 y =500
		fmt.Println("poolPrice > depositPrice")
		if xDeposit.Cmp(yDeposit) == -1 {
			rightY := y.Add(x, poolPrice)
			lpBalance := xDeposit.Sqrt(xDeposit.Mul(xDeposit, rightY))
			lp = &token{
				Name:    p.getPoolCoinName(),
				Balance: lpBalance,
			}
			p.Rs.X.Balance = x.Add(x, xDeposit)
			p.Rs.Y.Balance = y.Add(x, rightY)
			p.PC.Balance = pc.Add(pc, lpBalance)
		}
	case 0:
		fmt.Println("Equal!")
		//poolPrice == depositPrice
		p.Rs.X.Balance = x.Add(x, xDeposit)
		p.Rs.Y.Balance = y.Add(y, yDeposit)

		//lpBalance = sqrt(x*y)
		lpBalance := xDeposit.Sqrt(xDeposit.Mul(xDeposit, yDeposit))
		fmt.Println("lpBalance = ", lpBalance)
		p.PC.Balance = pc.Add(pc, lpBalance)
		lp = &token{
			Name:    p.getPoolCoinName(),
			Balance: lpBalance,
		}
	//poolPrice < depositPrice
	case 1:
		fmt.Println("poolPrice < depositPrice")
		//lpBalance = sqrt(x*y)
		var resultY *big.Float
		resultY = resultY.Mul(xDeposit, poolPrice)
		lpBalance := xDeposit.Sqrt(xDeposit.Mul(xDeposit, resultY))
		p.Rs.X.Balance = x.Add(x, xDeposit)
		p.Rs.Y.Balance = y.Add(y, resultY)
		p.PC.Balance = pc.Add(pc, lpBalance)
		lp = &token{
			Name:    p.getPoolCoinName(),
			Balance: lpBalance,
		}
	}
	fmt.Println("----------------------")
	return lp
}

func (p *Pool) WithDraw(lp token) (x, y token) {
	if lp.Name != p.getPoolCoinName() {
		log.Panic("not pool coin")
	}
	poolBalance := p.getPoolCoinBalance()
	rx, ry := p.Reserve()
	// (lp / poolBalance) *100
	// ex
	// pool.balance = 1000
	// lp = 10
	// percent = 1
	// (10/1000) * 100 = 1
	lpBalance := lp.GetTokenBalance()
	fmt.Println("lpBalance=", lpBalance)
	percent := getPercent(lpBalance, poolBalance)
	fmt.Println("percent =", percent)
	//xPrice = x*percent / 100
	//ex
	//x = 1000
	//percent = 1
	//(x * percent) / 100 = 10
	xPrice := getBalanceFromPercent(rx, percent)
	fmt.Println("xPrice = ", xPrice)
	yPrice := getBalanceFromPercent(ry, percent)
	fmt.Println("yPrice=", yPrice)
	xName, yName := p.getPairNameFromPool()
	x = token{
		Name:    xName,
		Balance: xPrice,
	}
	y = token{
		Name:    yName,
		Balance: yPrice,
	}
	p.Rs.X.Balance = rx.Sub(rx, xPrice)
	p.Rs.Y.Balance = ry.Sub(ry, yPrice)
	p.PC.Balance = poolBalance.Sub(poolBalance, lp.Balance)
	return x, y
}

func (p *Pool) Swap(t token) token {
	xName, yName := p.getPairNameFromPool()
	tName := t.Name
	if tName != xName || t.Name != yName {
		log.Panicf("your token is to exsits at pool")
	}

	//This is what you want to swap balance of token
	tBalance := t.GetTokenBalance()

	//x,y balance
	rx, ry := p.Reserve()

	//price
	k := p.K()
	if tName == xName {
		//+rx -> - ry
		rx = rx.Add(rx, tBalance)
		// sendY := ry.Sub(ry, k.Div(k, rx))
		sendY := ry.Sub(ry, k.Quo(k, rx))
		fee := sendY.Quo(sendY, big.NewFloat(0.03))
		sendY = sendY.Sub(sendY, fee)
		return token{
			Name:    yName,
			Balance: sendY,
		}
	} else {
		//+ry -> -rx
		ry = ry.Add(rx, tBalance)
		sendX := rx.Sub(rx, k.Quo(k, ry))
		fee := sendX.Quo(sendX, big.NewFloat(0.03))
		sendX = sendX.Sub(sendX, fee)
		return token{
			Name:    xName,
			Balance: sendX,
		}
	}
}
