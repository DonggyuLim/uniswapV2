package pool

import (
	"fmt"

	"cosmossdk.io/math"
	m "cosmossdk.io/math"
	u "github.com/DonggyuLim/uniswap/utils"
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
func (p *Pool) mintLP(account string, dx, dy m.Uint) {
	x, y := p.Reserve()
	lpTotal := p.LP.GetTotalSupply()
	if p.LP.GetTotalSupply().Equal(m.ZeroUint()) {
		// LP := uint64(u.Sqrt(dx, dy))
		LP := dx.Mul(dy)
		p.LP.Mint(account, LP)
		return
	}
	xQuo := (x.Quo(dx)).Mul(lpTotal)
	yQuo := (y.Quo(dy)).Mul(lpTotal)

	if xQuo.GT(yQuo) {
		p.LP.Mint(account, yQuo)
		return
	}
	p.LP.Mint(account, xQuo)
}

func (p *Pool) Deposit(account string, tokenX, tokenY Token) error {

	//현재 저장되어 있는 X,Y 토큰 수량
	rx, ry := p.Reserve()

	//deposit 하고 싶은 수량

	//!!!!!!!!!!!!
	//dx,dy 가 있는지 확인해야함.

	x, y := tokenX.GetBalance(), tokenY.GetBalance()
	dx, dy := math.LegacyMustNewDecFromStr(x.String()), math.LegacyMustNewDecFromStr(y.String())

	//balance check
	// account,amount -> grc20
	// amount > balance = err
	err := grc20TwoBalanceCheck(account, tokenX, tokenY)
	if err != nil {
		return err
	}

	//poolPirce
	pp := p.poolPrice()

	//deposit 한 토큰쌍의 가격
	dp := u.GetPrice(x, y)

	// 가격이 맞지 않으면 에러를 리턴해버리는 방법도 있을듯
	//그러나 자동으로 디파짓 할 수 있게 해주면 더 좋을듯

	// pp == dp 가 동일한 경우
	if pp.Equal(dp) {
		err := p.gprcDeposit(tokenX.GetTokenName(), account, x)
		if err != nil {
			return err
		}
		err = p.gprcDeposit(tokenY.GetTokenName(), account, y)
		// err = sendApprove(tokenY.GetTokenName(), account, p.GetName(), dy)
		if err != nil {
			return err
		}

		//풀 토큰 발행
		p.mintLP(account, x, y)
		p.X.Balance = rx.Add(x)
		p.Y.Balance = ry.Add(y)
		return nil
	}

	//pp 가 1보다 큰 decimal 일 경우
	if pp.GT(math.LegacyOneDec()) {
		y = m.Uint(dx.QuoTruncate(pp).TruncateInt())
		err = p.gprcDeposit(tokenX.GetTokenName(), account, x)
		if err != nil {
			return err
		}
		err = p.gprcDeposit(tokenY.GetTokenName(), account, y)
		// err = sendApprove(tokenY.GetTokenName(), account, p.GetName(), dy)
		if err != nil {
			return err
		}

		p.mintLP(account, x, y)
		p.X.Balance = rx.Add(x)
		p.Y.Balance = ry.Add(y)
		return nil
	}

	// pp 가 1보다 작은 decimal 일경우
	x = m.Uint(dy.Mul(pp).TruncateInt())
	err = p.gprcDeposit(tokenX.GetTokenName(), account, x)
	if err != nil {
		return err
	}
	err = p.gprcDeposit(tokenY.GetTokenName(), account, y)
	// err = sendApprove(tokenY.GetTokenName(), account, p.GetName(), dy)
	if err != nil {
		return err
	}
	p.mintLP(account, x, y)
	p.X.Balance = rx.Add(x)
	p.Y.Balance = ry.Add(y)
	return nil
}

// 원래 스왑이 일어나면 아마 approve 로 먼저 소유권 이전한 후에
// allowance 로 확인한 뒤에 후 처리를 할 것이다.
// 그런데 현재 내 프로그램에선 계정에서 approve 쏴주기가 불편하므로
// 이 함수안에서 다 처리하는 로직으로 만듬.
func (p *Pool) gprcDeposit(tokenName, account string, amount m.Uint) error {
	err := GRPCsendApprove(tokenName, account, p.GetName(), amount)
	if err != nil {
		return err
	}
	allowance, err := GRPCgetAllowance(tokenName, account, p.GetName())
	if allowance.LT(amount) || err != nil {
		return err
	}
	err = GRPCsendTransferFrom(tokenName, account, p.GetName(), p.GetName(), amount)
	if err != nil {
		return err
	}
	return nil
}

// WithDraw는 account가 가지고 있는 lp 비율만큼 리턴받음.
func (p *Pool) WithDraw(account string, amount math.Uint) error {

	//account balance

	err := p.lpCheckBalance(account, amount)
	if err != nil {
		return err
	}
	//pool totalsupply
	ps := p.LP.GetTotalSupply()
	rx, ry := p.Reserve()

	percent := u.GetPercent(amount, ps)

	//xb = will send X
	xb := u.GetBalanceFromPercent(rx, percent)
	yb := u.GetBalanceFromPercent(ry, percent)

	//grc20 데이터 베이스에 보내줘야함.
	err = GRPCsendTransfer(p.X.Name, p.GetName(), account, xb)
	if err != nil {
		return err
	}
	err = GRPCsendTransfer(p.Y.Name, p.GetName(), account, yb)
	if err != nil {
		return err
	}
	p.X.Balance = rx.Sub(xb)
	p.Y.Balance = ry.Sub(yb)
	p.LP.Burn(account, amount)
	return nil
}

// token = 교환할 토큰
func (p *Pool) Swap(tokenName, account string, amount math.Uint) error {
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

	//checkingbalance
	err := grc20CheckBalance(tokenName, account, amount)
	if err != nil {
		return err
	}
	//pool X,Y balance
	rx, ry := p.Reserve()

	// pool K
	k := m.Uint(p.K())

	xName, yName := p.getPairNameFromPool()
	//approve 를 제공 받고
	//allowance 확인을 하고
	//transferfrom 으로 balance 로 데이터 이동 시킨후
	//받고싶어하는 코인을 보내주면됨.
	if tokenName == xName {
		//+rx   -ry
		fmt.Println(tokenName)
		rx = rx.Add(amount)

		//

		sendY := ry.Sub(k.Quo(rx))
		fmt.Println("sendY=======", sendY)
		//fee 는 어디다 모아둘까?

		fee := sendY.Mul(u.GetBalanceFromPercent(sendY, p.FeeRate))
		fmt.Println("fee=======", fee)
		sendY = sendY.Sub(fee)
		err := p.grpcSwap(xName, yName, account, amount, sendY)
		if err != nil {
			return err
		}

		//fee 는 풀에서 데이터로 저장하고 잇다가 approve 해줘서 넘겨주면 될듯
		//x 토큰을 받았으니 y 토큰 피를 올려줘야함.
		p.YFee = p.XFee.Add(fee)
		p.X.Balance = rx
		p.Y.Balance = ry.Sub(sendY)
		return nil
	} else {

		//-rx   ry+
		ry = ry.Add(amount)
		sendX := rx.Sub(k.Quo(ry))

		fee := sendX.Mul(u.GetBalanceFromPercent(sendX, p.FeeRate))
		sendX = sendX.Sub(fee)
		err := p.grpcSwap(yName, xName, account, amount, sendX)
		if err != nil {
			return err
		}
		p.XFee = p.XFee.Add(fee)
		p.X.Balance = rx.Sub(sendX)
		p.Y.Balance = ry
		return nil
	}
}

// sendToken = 이용자가 보낸 토큰
// wantToken = 원하는 토큰
// sendAmount = 이용자가 보낸 토큰 양
// reciveAmount = 이용자가 받아야할 양
func (p *Pool) grpcSwap(sendTokenName, wantTokenName, account string, sendAmount, reciveAmount math.Uint) error {
	err := GRPCsendApprove(sendTokenName, account, p.GetName(), sendAmount)
	if err != nil {
		return err
	}
	allowance, err := GRPCgetAllowance(sendTokenName, account, p.GetName())
	if allowance.LT(sendAmount) || err != nil {
		return err
	}
	err = GRPCsendTransferFrom(sendTokenName, account, p.GetName(), p.GetName(), sendAmount)
	if err != nil {
		return err
	}
	err = GRPCsendTransfer(wantTokenName, p.GetName(), account, reciveAmount)
	if err != nil {
		return err
	}
	return nil
}
