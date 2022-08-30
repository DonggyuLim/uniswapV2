package pool

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"
)

const fee = "0.03"

// business logic
// 난중에DB 업데이트해야함.
func (p *Pool) Deposit(tokenA, tokenB token) (lp *token) {

	x, y := p.Reserve() //pool reserve

	xDeposit, yDeposit := tokenA.Balance, tokenB.Balance
	fmt.Printf("xDeposit = %v , yDeposit =%v\n", xDeposit, yDeposit)
	poolPrice, err := price(x, y)

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("poolPrice = %v\n", poolPrice)

	depositPrice, err := price(xDeposit, yDeposit)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("depositPrice = %v\n", depositPrice)

	pc := p.getPoolCoinBalance()

	switch poolPrice.Cmp(depositPrice) {
	//poolPrice < depositPrice
	case -1:
		// x = 30 y =500
		fmt.Println("poolPrice > depositPrice")
		if xDeposit.Cmp(yDeposit) == -1 {
			// rightY := y.Add(x, poolPrice)
			rightY := x.Add(poolPrice)
			lpBalance := sqrt(xDeposit.Mul(rightY))

			lp = &token{
				Name:    p.getPoolCoinName(),
				Balance: lpBalance,
			}
			p.Rs.X.Balance = x.Add(xDeposit)
			p.Rs.Y.Balance = y.Add(rightY)
			p.PC.Balance = pc.Add(lpBalance)
		}
	case 0:
		fmt.Println("Equal!")
		//poolPrice == depositPrice
		p.Rs.X.Balance = x.Add(xDeposit)
		p.Rs.Y.Balance = y.Add(yDeposit)

		//lpBalance = sqrt(x*y)
		lpBalance := sqrt(xDeposit.Mul(yDeposit))

		fmt.Println("lpBalance = ", lpBalance)
		p.PC.Balance = pc.Add(lpBalance)
		lp = &token{
			Name:    p.getPoolCoinName(),
			Balance: lpBalance,
		}
	//poolPrice < depositPrice
	case 1:
		fmt.Println("poolPrice < depositPrice")
		//lpBalance = sqrt(x*y)

		resultY := xDeposit.Mul(poolPrice)

		lpBalance := sqrt(xDeposit.Mul(resultY))
		p.Rs.X.Balance = x.Add(xDeposit)
		p.Rs.Y.Balance = y.Add(resultY)
		p.PC.Balance = pc.Add(lpBalance)
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
	p.Rs.X.Balance = rx.Sub(xPrice)
	p.Rs.Y.Balance = ry.Sub(yPrice)
	p.PC.Balance = poolBalance.Sub(lpBalance)
	return x, y
}

func (p *Pool) Swap(t token) token {
	xName, yName := p.getPairNameFromPool()
	tName := t.Name
	if tName != xName || t.Name != yName {
		log.Panicf("your token is to exsits at pool")
	}

	//This is what you want to swap balance of token
	tB := t.GetTokenBalance()

	//x,y balance
	rx, ry := p.Reserve()

	//price
	k := p.K()
	if tName == xName {
		//+rx -> - ry
		rx = rx.Add(tB)
		// sendY := ry.Sub(ry, k.Div(k, rx))
		sendY := ry.Sub(k.Div(rx))
		fee, err := decimal.NewFromString(fee)
		if err != nil {
			log.Panic(err)
		}
		fee = sendY.Div(fee)
		sendY = sendY.Sub(fee)
		return token{
			Name:    yName,
			Balance: sendY,
		}
	} else {
		//+ry -> -rx
		ry = ry.Add(tB)
		sendX := rx.Sub(k.Div(ry))
		fee, err := decimal.NewFromString(fee)
		if err != nil {
			log.Panic(err)
		}
		fee = sendX.Div(fee)

		sendX = sendX.Sub(fee)
		return token{
			Name:    xName,
			Balance: sendX,
		}
	}
}
