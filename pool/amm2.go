package pool

import (
	"errors"

	"github.com/DonggyuLim/uniswap/math"
	"github.com/shopspring/decimal"
)

//Reference
// https://uniswap.org/whitepaper.pdf
// - uniswap v2 whitepaper 3.4 Initialization of liquidity token supply

/*
    참고 - uniswap v2
    https://docs.uniswap.org/protocol/V2/concepts/core-concepts/pools
   유동성이 풀에 예치될 때마다 유동성 토큰으로 알려진 고유한 토큰 이 발행되어 공급자의 주소로 전송됩니다. 이 토큰은 풀에 대한 주어진 유동성 공급자의 기여도를 나타냅니다.
    제공된 풀의 유동성 비율은 공급자가 받는 유동성 토큰의 수를 결정합니다.
   공급자가 새 풀을 발행하는 경우 받을 유동성 토큰의 수는 sqrt(X * Y)와 같으며 여기서 x와 y는 제공된 각 토큰의 양을 나타냅니다.

   //새 풀이 아닐 경우
   // (x 넣은 수량 / x 전체 수량) * LP 토큰 총량
*/

//empty pool 일 경우 처음 공급자는
//√X*Y

//empty 풀이 아닐 경우
// (x 토큰 넣기전 수량 / 넣을 x 토큰 수량) * 풀에 있는 토큰 량

// 위의 공식은 x 가 될수도 y 가 될 수도 있음.
// 더 작은 쪽으로 LP 토큰 제공.
func (p *Pool) mintLP(address string, dx, dy decimal.Decimal) {

	one := decimal.NewFromInt(0)
	if p.LP.GetTotalSupply().Cmp(one) != 1 {

		LP := math.Sqrt(dx.Mul(dy))
		p.LP.Mint(address, LP)
		return
	} else {

		caseX := dx.Div(p.X.GetBalance()).Mul(p.LP.GetTotalSupply())
		caseY := dy.Div(p.Y.GetBalance()).Mul(p.LP.GetTotalSupply())
		switch caseX.Cmp(caseY) {
		case -1:
			p.LP.Mint(address, caseX)
		case 0, 1:
			p.LP.Mint(address, caseY)
		}
	}
}

func (p *Pool) Deposit(address string, tokenA, tokenB Token) error {

	//현재 저장되어 있는 X,Y 토큰 수량
	rx, ry := p.Reserve()

	//deposit 하고 싶은 수량

	//!!!!!!!!!!!!
	//dx,dy 가 있는지 확인해야함.

	dx, dy := tokenA.GetBalance(), tokenB.GetBalance()

	//balance check
	// address,amount -> grc20
	// amount > balance = err
	err := depositBalanceCheck(address, tokenA, tokenB)
	if err != nil {
		return err
	}

	//poolPirce
	pp := p.poolPrice()

	//deposit 한 토큰쌍의 가격
	dp := math.Price(dx, dy)

	// 가격이 맞지 않으면 에러를 리턴해버리는 방법도 있을듯
	//그러나 자동으로 디파짓 할 수 있게 해주면 더 좋을듯

	switch pp.Cmp(dp) {
	case 0:
		// pp == dp

		//grpc에 approve 호출해야함.
		err = sendApprove(tokenA.TokenName(), address, dx)
		if err != nil {
			return err
		}
		err = sendApprove(tokenB.TokenName(), address, dy)
		if err != nil {
			return err
		}

		//풀 토큰 발행
		p.mintLP(address, dx, dy)
		p.X.Balance = rx.Add(dx)
		p.Y.Balance = ry.Add(dy)
		return nil
	case 1:
		// pp 가 1이상이면 나눠줘야하고 1미만이면 곱해줘야 비율이 맞춰짐.
		var tempY decimal.Decimal
		one := decimal.NewFromInt(1)
		if pp.Cmp(one) == 1 && dp.Cmp(one) == -1 {
			tempY = dx.Mul(pp)
		} else {
			tempY = dx.Div(pp)
		}
		//만약 dy - tempY 가 음수라면 패닉일으켜야해.
		//왜냐하면 필요한 Y 가 더 크다는 뜻이니까.
		if dy.Sub(tempY).Cmp(one) == -1 {
			err = errors.New("send more token b")
			return err
		}

		//grpc에 approve 호출해야함.

		err = sendApprove(tokenA.TokenName(), address, dx)
		if err != nil {
			return err
		}
		err = sendApprove(tokenB.TokenName(), address, tempY)
		if err != nil {
			return err
		}
		p.mintLP(address, dx, tempY)
		p.X.Balance = rx.Add(dx)
		p.Y.Balance = ry.Add(tempY)
	case -1:
		tempX := dx.Mul(pp)
		one := decimal.NewFromInt(1)
		if dx.Sub(tempX).Cmp(one) == -1 {

			err = sendApprove(tokenA.TokenName(), address, tempX)
			if err != nil {
				return err
			}
			err = sendApprove(tokenB.TokenName(), address, dy)
			if err != nil {
				return err
			}
			p.mintLP(address, tempX, dy)
			p.X.Balance = rx.Add(tempX)
			p.Y.Balance = ry.Add(dy)
		}
	}
	return nil
}

func (p *Pool) WithDraw(address string, amount decimal.Decimal) error {

	//address balance
	ab := p.LP.BalanceOf(address)
	err := lpCheckBalance(ab, amount)
	if err != nil {
		return err
	}
	//pool totalsupply
	ps := p.LP.GetTotalSupply()
	rx, ry := p.Reserve()

	percent := math.GetPercent(amount, ps)

	//xb = will send X
	xb := math.GetBalanceFromPercent(rx, percent)
	yb := math.GetBalanceFromPercent(ry, percent)

	//grc20 데이터 베이스에 보내줘야함.
	err = sendTransferFrom(p.X.Name, p.GetName(), address, p.GetName(), xb)
	if err != nil {
		return err
	}
	err = sendTransferFrom(p.Y.Name, p.GetName(), address, p.GetName(), yb)
	if err != nil {
		return err
	}
	p.X.Balance = rx.Sub(xb)
	p.Y.Balance = ry.Sub(yb)
	p.LP.Burn(address, amount)
	return nil
}

// token = 교환할 토큰
func (p *Pool) Swap(token, address string, amount decimal.Decimal) error {
	//스왑가격 결정 방법
	//X = 10 Y = 100 k = 1000
	// 한개의 X 를 보내면 얼마만큼의 Y 를 받을 수 있을까?
	// X 는 11이 될 것이다. Y 는 몇개가 되어야하지?
	// uniswap v2 에서는 항상 k 는 일정해야한다.
	// 11x * ?Y = 1000 이 되어야하는 것이다.
	//그렇다면 k/?Y = 11x 가 된다.
	// 1000을 11로 나누면? 90.9090909091
	//중요한점은 토큰 을 제공하고 나서 90.9090909091 가 되어야 한다는 것이다.
	//그러니까 원래 Y 에서 90.9090909091 가 되는 토큰 수량 만큼 보내주면된다.
	// 그리고 여기서 0.3 %를 때고 제공을 해준다.

	//pool X,Y balance
	rx, ry := p.Reserve()

	//Price
	k := p.K()

	xName, yName := p.getPairNameFromPool()
	if token == xName {
		//+rx   -ry
		rx = rx.Add(amount)
		sendY := ry.Sub(k.Div(rx))
		//fee 는 어디다 모아둘까?
		fee := sendY.Mul(p.Fee)
		p.LP.Balance["0"] = p.LP.Balance["0"].Add(fee)
		sendY = sendY.Sub(fee)
		err := sendTransferFrom(yName, p.GetName(), address, p.GetName(), sendY)

		if err != nil {
			return err
		}
		p.X.Balance = rx
		p.Y.Balance = ry.Sub(sendY)

		return nil
	} else {
		//-rx   ry+
		ry = ry.Add(amount)
		sendX := rx.Sub(k.Div(ry))
		fee := sendX.Mul(p.Fee)
		sendX = sendX.Sub(fee)
		err := sendTransferFrom(xName, p.GetName(), address, p.GetName(), sendX)
		if err != nil {
			return err
		}
		p.X.Balance = rx.Sub(sendX)
		p.Y.Balance = ry
		return nil
	}
}
