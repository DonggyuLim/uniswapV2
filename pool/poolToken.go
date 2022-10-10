package pool

import (
	"errors"

	m "cosmossdk.io/math"
	u "github.com/DonggyuLim/uniswap/utils"
)

type poolToken struct {
	Name        string
	Symbol      string
	Decimal     uint8
	TotalSupply m.Uint
	Balance     map[string]m.Uint
	Allowances  map[string]m.Uint
}

func NewPoolToken(name, symbol string, dec uint8) poolToken {
	bm := make(map[string]m.Uint)
	am := make(map[string]m.Uint)
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
			err := GRPCsendApprove(tokenAName, k, poolName, offerBalance)
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
func (t *poolToken) GetTotalSupply() m.Uint {
	return t.TotalSupply
}

func (t *poolToken) BalanceOf(account string) m.Uint {
	return t.Balance[account]
}

func (t *poolToken) GetDecimal() uint8 {
	return t.Decimal
}

func (t *poolToken) Allowance(owner, spender string) m.Uint {
	return t.allowance(owner, spender)
}
func (t *poolToken) allowance(owner, spender string) m.Uint {
	key := owner + ":" + spender
	return t.Allowances[key]
}

// /////////////////////////////////////////////////
// Mutate

func (t *poolToken) Transfer(from, to string, amount m.Uint) error {
	err := t.checkBalance(from, amount)
	if err != nil {
		return err
	}
	t.transfer(from, to, amount)
	return nil
}

func (t *poolToken) transfer(from, to string, amount m.Uint) {

	fromBalance := t.Balance[from]
	t.Balance[from] = fromBalance.Sub(amount)
	toBalance := t.Balance[to]
	t.Balance[to] = toBalance.Add(amount)
}

func (t *poolToken) Approve(owner, spender string, amount m.Uint) error {
	if err := t.checkBalance(owner, amount); err != nil {
		return err
	}
	t.approve(owner, spender, amount)

	return nil
}

func (t *poolToken) approve(owner, spender string, amount m.Uint) error {

	key := owner + ":" + spender
	currentBalance := t.Allowances[key]
	t.Balance[owner] = t.BalanceOf(owner).Sub(amount)
	t.Allowances[key] = currentBalance.Add(amount)
	return nil
}

func (t *poolToken) TransferFrom(owner, spender, to string, amount m.Uint) error {
	if err := t.checkspendAllowance(owner, spender, amount); err != nil {
		return err
	}
	t.transferfrom(owner, spender, to, amount)
	return nil
}

func (t *poolToken) transferfrom(owner, spender, to string, amount m.Uint) {
	key := owner + ":" + spender
	t.Allowances[key] = t.allowance(owner, spender).Sub(amount)

	if t.Allowances[key].LT(m.NewUint(0)) {
		delete(t.Allowances, key)
	}
	t.Balance[spender] = t.BalanceOf(spender).Add(amount)
}

func (t *poolToken) Mint(account string, amount m.Uint) {
	t.mint(account, amount)
}

func (t *poolToken) mint(address string, amount m.Uint) {

	t.TotalSupply = t.TotalSupply.Add(amount)
	t.Balance[address] = t.BalanceOf(address).Add(amount)
}

func (t *poolToken) Burn(address string, amount m.Uint) {
	t.burn(address, amount)
}

func (t *poolToken) burn(address string, amount m.Uint) {
	currentBalance := t.BalanceOf(address)
	newBalance := currentBalance.Sub(amount)

	if newBalance.LT(m.NewUint(0)) {
		t.TotalSupply = t.GetTotalSupply().Sub(amount)
		t.Balance[address] = m.NewUint(0)
	} else {
		t.TotalSupply = t.GetTotalSupply().Sub(amount)
		t.Balance[address] = m.Uint(newBalance)
	}
}

// //////////////////////////////////////
func (t *poolToken) checkBalance(owner string, amount m.Uint) error {
	balance := t.BalanceOf(owner)
	if balance.LT(amount) {
		return errors.New("amount is biger than owner balance")
	}
	return nil
}

func (t *poolToken) checkspendAllowance(owner, spender string, amount m.Uint) error {
	allowance := t.allowance(owner, spender)
	if allowance.LT(amount) {
		return errors.New("amount is biger than allowance")
	}
	return nil
}
