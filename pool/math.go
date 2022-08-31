package pool

import (
	"github.com/shopspring/decimal"
)

func price(x, y decimal.Decimal) (price decimal.Decimal) {

	price = x.Div(y)
	return price
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
	// xF := x.InexactFloat64()
	xF := x.BigFloat()
	sqrt, _ := xF.Sqrt(xF).Float64()
	// sqrt := math.Sqrt(xF)
	return decimal.NewFromFloat(sqrt)
}
