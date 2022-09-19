package pool

import (
	"errors"

	u "github.com/DonggyuLim/uniswap/utils"
	"github.com/shopspring/decimal"
)

type poolToken struct {
	Name        string
	Symbol      string
	Decimal     uint8
	TotalSupply decimal.Decimal
	Balance     map[string]decimal.Decimal
	Allowances  map[string]decimal.Decimal
}

func NewPoolToken(name, symbol string, dec uint8) poolToken {
	bm := make(map[string]decimal.Decimal)
	am := make(map[string]decimal.Decimal)
	return poolToken{
		Name:       name,
		Symbol:     symbol,
		Decimal:    dec,
		Balance:    bm,
		Allowances: am,
	}
}

func (t *poolToken) ShareFee(tokenAName, tokenBName, poolName string) {
	feeBalance := t.BalanceOf("0x")
	for k, v := range t.Balance {
		//
		if k == "0x" {
			continue
		} else {
			percent := u.GetPercent(v, t.TotalSupply)
			offerBalance := u.GetBalanceFromPercent(feeBalance, percent)
			err := sendApprove(tokenAName, k, poolName, offerBalance)
			if err != nil {
				panic(err)
			}

		}
	}
}

// /////////////////
// ////////////////
// ////////////////
// Query
func (t *poolToken) GetName() string {
	return t.Name
}
func (t *poolToken) GetSymbol() string {
	return t.Symbol
}
func (t *poolToken) GetTotalSupply() decimal.Decimal {
	return t.TotalSupply
}

func (t *poolToken) BalanceOf(account string) decimal.Decimal {
	return t.Balance[account]
}

func (t *poolToken) GetDecimal() uint8 {
	return t.Decimal
}

func (t *poolToken) Allowance(owner, spender string) decimal.Decimal {
	return t.allowance(owner, spender)
}
func (t *poolToken) allowance(owner, spender string) decimal.Decimal {
	key := owner + ":" + spender
	return t.Allowances[key]
}

// /////////////////////////////////////////////////
// Mutate

func (t *poolToken) Transfer(from, to string, amount decimal.Decimal) error {
	err := t.checkBalance(from, amount)
	if err != nil {
		return err
	}
	t.transfer(from, to, amount)
	return nil
}

func (t *poolToken) transfer(from, to string, amount decimal.Decimal) {

	fromBalance := t.Balance[from]
	t.Balance[from] = fromBalance.Sub(amount)
	toBalance := t.Balance[to]
	t.Balance[to] = toBalance.Add(amount)
}

func (t *poolToken) Approve(owner, spender string, amount decimal.Decimal) error {
	if err := t.checkBalance(owner, amount); err != nil {
		return err
	}
	t.approve(owner, spender, amount)

	return nil
}

func (t *poolToken) approve(owner, spender string, amount decimal.Decimal) error {

	key := owner + ":" + spender
	currentBalance := t.Allowances[key]
	t.Balance[owner] = t.BalanceOf(owner).Sub(amount)
	t.Allowances[key] = currentBalance.Add(amount)
	return nil
}

func (t *poolToken) TransferFrom(owner, spender, to string, amount decimal.Decimal) error {
	if err := t.checkspendAllowance(owner, spender, amount); err != nil {
		return err
	}
	t.transferfrom(owner, spender, to, amount)
	return nil
}

func (t *poolToken) transferfrom(owner, spender, to string, amount decimal.Decimal) {
	key := owner + ":" + spender
	t.Balance[to] = t.BalanceOf(spender).Add(amount)
	t.Allowances[key] = t.allowance(owner, spender).Sub(amount)
	zero := decimal.NewFromInt(0)
	if t.Allowances[key].Cmp(zero) != 1 {
		delete(t.Allowances, key)
	}
}

func (t *poolToken) Mint(account string, amount decimal.Decimal) {

	t.mint(account, amount)
}

func (t *poolToken) mint(address string, amount decimal.Decimal) {

	t.TotalSupply = t.TotalSupply.Add(amount)
	t.Balance[address] = t.BalanceOf(address).Add(amount)
}

func (t *poolToken) Burn(address string, amount decimal.Decimal) {
	t.burn(address, amount)
}

func (t *poolToken) burn(address string, amount decimal.Decimal) {
	currentBalance := t.BalanceOf(address)
	newBalance := currentBalance.Sub(amount)
	t.TotalSupply = t.GetTotalSupply().Sub(amount)
	t.Balance[address] = newBalance
}

// //////////////////////////////////////
func (t *poolToken) checkBalance(owner string, amount decimal.Decimal) error {
	balance := t.BalanceOf(owner)
	if balance.Cmp(amount) == -1 {
		return errors.New("amount is biger than owner balance")
	}
	return nil
}

func (t *poolToken) checkspendAllowance(owner, spender string, amount decimal.Decimal) error {
	Allowance := t.allowance(owner, spender)
	if Allowance.Cmp(amount) == -1 {
		return errors.New("amount is biger than allowance")
	}
	return nil
}
