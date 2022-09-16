package pool

import (
	"errors"

	"github.com/DonggyuLim/uniswap/client"
	u "github.com/DonggyuLim/uniswap/utils"
	"github.com/shopspring/decimal"
)

func depositBalanceCheck(address string, tokenA, tokenB Token) error {
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
	client := client.GetRPCClient()
	balance, err := client.GetBalance(tokenName, address)
	if err != nil {
		return err
	}
	if balance.Cmp(amount) == -1 {
		return err
	}
	return nil
}

func lpCheckBalance(balance, amount decimal.Decimal) error {
	if balance.Cmp(amount) == -1 {
		return errors.New("you have not enough lp")
	}
	return nil
}

func sendApprove(tokenName, account string, amount decimal.Decimal) error {
	_, err := client.GetRPCClient().Approve(tokenName, account, "0xpool", u.DecimalToUint64(amount))
	if err != nil {
		return err
	}
	return nil
}

func sendTransferFrom(tokenName, from, to, spender string, amount decimal.Decimal) error {
	_, err := client.GetRPCClient().TransferFrom(tokenName, from, to, spender, u.DecimalToUint64(amount))
	if err != nil {
		return err
	}
	return nil
}
