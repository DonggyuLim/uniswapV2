package pool

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/DonggyuLim/uniswap/client"
)

func GRPCsendApprove(tokenName, owner, spender string, amount math.Uint) error {
	fmt.Println("Approve!")
	_, err := client.GetClient().Approve(tokenName, owner, spender, amount.String())
	if err != nil {
		return err
	}
	return nil
}

func GRPCgetAllowance(tokenName, account, poolName string) (math.Uint, error) {
	fmt.Println("Allowance!")
	balance, err := client.GetClient().Allowance(tokenName, account, poolName)
	if err != nil {
		return math.NewUint(0), err
	}

	return math.NewUintFromString(balance), err
}

func GRPCsendTransferFrom(tokenName, owner, spender, to string, amount math.Uint) error {
	fmt.Println("TransferFrom!")
	_, err := client.GetClient().TransferFrom(tokenName, owner, spender, to, amount.String())
	if err != nil {
		return err
	}
	return nil
}

func GRPCsendTransfer(tokenName, from, to string, amount math.Uint) error {
	fmt.Println("Transfer!")
	_, err := client.GetClient().Transfer(tokenName, from, to, amount.String())
	if err != nil {
		return err
	}
	return nil
}
