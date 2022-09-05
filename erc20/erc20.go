package erc20

import "github.com/shopspring/decimal"

type address string

type Token struct {
	name    string
	symbol  string
	balance map[address]decimal.Decimal
	// allowance    map[address]map[address]decimal.Decimal
	totalBalance decimal.Decimal
}

func Initalize(name, symbol, account string, total int64) *Token {
	//balance map
	bm := make(map[address]decimal.Decimal)

	token := &Token{
		name,
		symbol,
		bm,
		decimal.NewFromInt(int64(total)),
	}
	token.mint(account, total)
	return token
}

func (t *Token) Name() string {
	return t.name
}

func (t *Token) Symbol() string {
	return t.symbol
}

func (t *Token) Decimal() int64 {
	return 10
}

func (t *Token) TotalSupply() int64 {
	return t.totalBalance.IntPart()
}

func (t *Token) BalanceOf(addr address) decimal.Decimal {
	// db.Get()
	return t.balance[addr]
}

func (t *Token) Transfer(to, from address, amount decimal.Decimal) {
	t.balance[to] = t.balance[to].Sub(amount)
	t.balance[from] = t.balance[from].Add(amount)
}

func (t *Token) mint(addr string, amount int64) {
	t.balance[address(addr)] = decimal.NewFromInt(amount)
}

// func (t *Token) Allowance(owner, spender address) {

// }

// func (t *Token) Approve(spender address, amount decimal.Decimal) {

// }

// func (t *Token) TransferFrom(from, to address, amount decimal.Decimal) {

// }

// func (t *Token) IncreaseAllowance(spender address, addedValue decimal.Decimal) {

// }
// func (t *Token) DecreaseAllowance(spender address, subtractedValue decimal.Decimal) {

// }

// func (t *Token) burn(account address, amount decimal.Decimal) {

// }

// func (t *Token) approve(owner, spender address) {
// 	if t.balance[owner].Sign() == 0{
// 		panic()
// 	}
// }

// func spendAllowance(owner, spender address, amount decimal.Decimal) {

// }
