package pool

import (
	"log"

	"github.com/DonggyuLim/uniswap/math"
	"github.com/shopspring/decimal"
)

const fee = "0.03"

// business logic
// 난중에DB 업데이트해야함.
func (p *Pool) Deposit(tokenA, tokenB token) (lp token) {
	//reserve token amount
	rx, ry := p.Reserve() //pool reserve

	//deposit token amount
	dx, dy := tokenA.Balance, tokenB.Balance

	//pp = pool price poolPrice = 고정
	pp := p.poolPrice()

	//dp = Deposit Price
	dp := math.Price(dx, dy)

	pc := p.getPoolCoinBalance()

	switch pp.Cmp(dp) {
	case 0:
		//pp == dp
		p.Rs.X.Balance = rx.Add(dx)
		p.Rs.Y.Balance = ry.Add(dy)
		//rp = return lp
		rp := math.Sqrt(dx.Mul(dy))
		p.PC.Balance = pc.Add(rp)
		lp = token{
			Name:    p.getPoolCoinName(),
			Balance: rp,
		}
		return
		//pp > dp
	case 1:
		// pp 가 1이상이면 나눠줘야하고 1미만이면 곱해줘야지 비율이 맞춰짐.
		var tempY decimal.Decimal
		one := decimal.NewFromInt(1)
		if pp.Cmp(one) == 1 && dp.Cmp(one) == -1 {
			tempY = dx.Mul(pp)
		} else {
			tempY = dx.Div(pp)
		}

		//만약 dy - tempY 가 음수라면 패닉일으켜야해.
		//왜냐하면 필요한 y 가 더 크다는 뜻이니까.
		if dy.Sub(tempY).Cmp(decimal.NewFromInt(0)) == -1 {
			log.Panic("Send more token B")
		}
		// dy - tempY 는 환불해줘야함.
		p.Rs.X.Balance = rx.Add(dx)
		p.Rs.Y.Balance = ry.Add(tempY)
		rp := math.Sqrt(dx.Mul(tempY))
		p.PC.Balance = pc.Add(rp)
		lp = token{
			Name:    p.getPoolCoinName(),
			Balance: rp,
		}
		return
		// pp < dp
	case -1:
		tempX := dy.Mul(pp)
		if dx.Sub(tempX).Cmp(decimal.NewFromInt(0)) == -1 {
			log.Panic("Send more token A")
		}
		p.Rs.X.Balance = rx.Add(tempX)
		p.Rs.Y.Balance = ry.Add(dy)
		rp := math.Sqrt(dy.Mul(tempX))
		lp = token{
			Name:    p.getPoolCoinName(),
			Balance: rp,
		}
		return
	}
	return
}

func (p *Pool) WithDraw(lp token) (x, y token) {
	if lp.Name != p.getPoolCoinName() {
		log.Panic("not pool coin")
	}
	//pb = pool.Balance
	pb := p.getPoolCoinBalance()
	rx, ry := p.Reserve()
	// (lp / pb) *100
	// ex
	// pool.balance = 1000
	// lp = 10
	// percent = 1
	// (10/1000) * 100 = 1
	//lb = send token(lp) balance
	lb := lp.TokenBalance()

	percent := math.GetPercent(lb, pb)

	//xPrice = x*percent / 100
	//ex
	//x = 1000
	//percent = 1
	//(x * percent) / 100 = 10

	//xb = will send x
	xb := math.GetBalanceFromPercent(rx, percent)

	//yb =will send y
	yb := math.GetBalanceFromPercent(ry, percent)

	xName, yName := p.getPairNameFromPool()
	x = token{
		Name:    xName,
		Balance: xb,
	}
	y = token{
		Name:    yName,
		Balance: yb,
	}
	p.Rs.X.Balance = rx.Sub(xb)
	p.Rs.Y.Balance = ry.Sub(yb)
	p.PC.Balance = pb.Sub(lb)
	return
}

// fee 는 따로 데이터베이스 업데이트 해줘야함.
func (p *Pool) Swap(t token) token {
	xName, yName := p.getPairNameFromPool()
	tName := t.Name
	if tName != xName || t.Name != yName {
		log.Panicf("your token is to exsits at pool")
	}

	//This is what you want to swap balance of token
	tB := t.TokenBalance()

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
