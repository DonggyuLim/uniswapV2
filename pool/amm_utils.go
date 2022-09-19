package pool

import (
	"errors"
	"fmt"

	"github.com/DonggyuLim/uniswap/client"
	u "github.com/DonggyuLim/uniswap/utils"
	"github.com/shopspring/decimal"
)

func grc20TwoBalanceCheck(address string, tokenA, tokenB Token) error {
	err := grc20CheckBalance(tokenA.Name, address, tokenA.GetBalance())
	if err != nil {
		return err
	}
	err = grc20CheckBalance(tokenB.Name, address, tokenB.GetBalance())
	if err != nil {
		return err
	}

	return nil
}

func grc20CheckBalance(tokenName, address string, amount decimal.Decimal) error {
	cli := client.GetClient()
	balance, err := cli.GetBalance(tokenName, address)

	if err != nil {
		return err
	}
	if balance.Cmp(amount) == -1 {
		return errors.New("balance isn't enough")
	}
	return nil
}

// Pool 의 메서드로 변경하기!
func lpCheckBalance(balance, amount decimal.Decimal) error {
	if balance.Cmp(amount) == -1 {
		return errors.New("you have not enough lp")
	}
	return nil
}

// sendToken = 이용자가 보낸 토큰
// wantToken = 원하는 토큰
// sendAmount = 이용자가 보낸 토큰 양
// reciveAmount = 이용자가 받아야할 양
func (p *Pool) swap(sendTokenName, wantTokenName, account string, sendAmount, reciveAmount decimal.Decimal) error {
	err := sendApprove(sendTokenName, account, p.GetName(), sendAmount)
	if err != nil {
		return err
	}
	allowance, err := getAllowance(sendTokenName, account, p.GetName())
	if allowance.Cmp(sendAmount) != 0 || err != nil {
		return err
	}
	err = sendTransferFrom(sendTokenName, account, p.GetName(), p.GetName(), sendAmount)
	if err != nil {
		return err
	}
	err = sendTransfer(wantTokenName, p.GetName(), account, reciveAmount)
	if err != nil {
		return err
	}
	return nil
}

// 원래 스왑이 일어나면 아마 approve 로 먼저 소유권 이전한 후에
// allowance 로 확인한 뒤에 후 처리를 할 것이다.
// 그런데 현재 내 프로그램에선 계정에서 approve 쏴주기가 불편하므로
// 이 함수안에서 다 처리하는 로직으로 만듬.
func (p *Pool) desposit(tokenName, account string, amount decimal.Decimal) error {
	err := sendApprove(tokenName, account, p.GetName(), amount)
	if err != nil {
		return err
	}
	allowance, err := getAllowance(tokenName, account, p.GetName())
	if allowance.Cmp(amount) != 0 || err != nil {
		return err
	}
	err = sendTransferFrom(tokenName, account, p.GetName(), p.GetName(), amount)
	if err != nil {
		return err
	}
	return nil
}

func sendApprove(tokenName, owner, spender string, amount decimal.Decimal) error {
	//p.name 으로로??
	fmt.Println("Approve!")
	_, err := client.GetClient().Approve(tokenName, owner, spender, u.DecimalToUint64(amount))
	if err != nil {
		return err
	}
	return nil
}

func getAllowance(tokenName, account, poolName string) (decimal.Decimal, error) {
	fmt.Println("Allowance!")
	balance, err := client.GetClient().Allowance(tokenName, account, poolName)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return balance, err
}

func sendTransferFrom(tokenName, owner, spender, to string, amount decimal.Decimal) error {
	fmt.Println("TransferFrom!")
	_, err := client.GetClient().TransferFrom(tokenName, owner, spender, to, u.DecimalToUint64(amount))
	if err != nil {
		return err
	}
	return nil
}

func sendTransfer(tokenName, from, to string, amount decimal.Decimal) error {
	fmt.Println("Transfer!")
	_, err := client.GetClient().Transfer(tokenName, from, to, u.DecimalToUint64(amount))
	if err != nil {
		return err
	}
	return nil
}
