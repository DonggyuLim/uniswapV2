package pool

import (
	"errors"

	m "cosmossdk.io/math"
	"github.com/DonggyuLim/uniswap/client"
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

func grc20CheckBalance(tokenName, address string, amount m.Uint) error {
	cli := client.GetClient()
	balancestr, err := cli.GetBalance(tokenName, address)
	balance := m.NewUintFromString(balancestr)
	if err != nil {
		return err
	}
	if balance.LT(amount) {
		return errors.New("balance isn't enough")
	}
	return nil
}
