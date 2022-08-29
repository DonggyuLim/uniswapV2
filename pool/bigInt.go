package pool

import (
	"math/big"
)

func Price(x, y int64) (*big.Int, error) {
	var price int64
	var err error
	if x > y {
		price = x / y
	} else if x == y {
		price = 1
	} else {
		price = y / x
	}

	return big.NewInt(price), err

}

// (x /tv ) * 100 = percent
// tv = tatal value
func getPercent(x, tv int64) *big.Int {
	floatx := float64(x)
	floattv := float64(tv)
	return big.NewInt(int64(floatx / floattv * 100.0))
}

// (tv * percentage) / 100
func getBalanceFromPercent(tv, percent int64) *big.Int {
	return big.NewInt((tv * percent) / 100)
}
