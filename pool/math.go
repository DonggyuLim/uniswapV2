package pool

import (
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

func price(x, y decimal.Decimal) (decimal.Decimal, error) {

	var price decimal.Decimal
	var err error

	switch x.Cmp(y) {
	case 1:
		price = x.Div(y)
		err = nil
	case 0:
		price = decimal.NewFromInt(1)
		err = nil
	case -1:
		err = errors.New("x must biger than y")
	}
	return price, err
}

// (x /tv ) * 100 = percent
// tv = tatal value
func getPercent(x, tv decimal.Decimal) decimal.Decimal {
	return x.Div(tv).Mul(decimal.NewFromInt(100))

}

// (tv * percentage) / 100
func getBalanceFromPercent(tv, percent decimal.Decimal) decimal.Decimal {
	return tv.Mul(percent).Div(decimal.NewFromInt(100))

}

func sqrt(x decimal.Decimal) decimal.Decimal {
	xF := x.InexactFloat64()
	sqrt := math.Sqrt(xF)
	return decimal.NewFromFloat(sqrt)
}
