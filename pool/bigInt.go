package pool

import (
	"math/big"
)

func Price(x, y *big.Float) (*big.Float, error) {
	fx, _ := x.Float64()
	fy, _ := y.Float64()
	var price float64
	var err error
	if fx > fy {
		price = fx / fy
	} else if fx == fy {
		price = 1
	} else {
		price = fy / fx
	}

	return big.NewFloat(price), err

}

// (x /tv ) * 100 = percent
// tv = tatal value
func getPercent(x, tv *big.Float) *big.Float {
	ftx, _ := x.Float64()
	ftv, _ := tv.Float64()
	return big.NewFloat(ftx / ftv * 100.0)
}

// (tv * percentage) / 100
func getBalanceFromPercent(tv, percent *big.Float) *big.Float {
	ft, _ := tv.Float64()
	fp, _ := percent.Float64()
	return big.NewFloat((ft * fp) / 100)

}
