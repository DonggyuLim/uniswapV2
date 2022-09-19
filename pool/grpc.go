package pool

import (
	"fmt"

	"github.com/DonggyuLim/uniswap/client"
	u "github.com/DonggyuLim/uniswap/utils"
	"github.com/shopspring/decimal"
)

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
