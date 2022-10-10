package utils

import (
	m "cosmossdk.io/math"
)

var decimalFromStr = m.LegacyMustNewDecFromStr

func GetPrice(x, y m.Uint) (price m.LegacyDec) {
	dx := decimalFromStr(x.String())
	dy := decimalFromStr(y.String())
	price = dx.Quo(dy)
	return
}

// (x /tv ) * 100 = percent
// tv = tatal value
func GetPercent(x, tv m.Uint) m.LegacyDec {
	dx := decimalFromStr(x.String())
	dtv := decimalFromStr(tv.String())

	return dx.Quo(dtv).MulInt64(100)
}

// (tv * percentage) / 100
func GetBalanceFromPercent(tv m.Uint, percent m.LegacyDec) m.Uint {
	dtv := decimalFromStr(tv.String())
	result := dtv.Mul(percent).QuoInt64(100).TruncateInt()
	// result := DecimalToUint64(tv.Mul(percent).Div(hundread))

	return ChanegeUintFromInt(result)
}

// sqrt return not exact value.....
func Sqrt(x, y m.Uint) (sqrt m.LegacyDec, err error) {
	mul := x.Mul(y).String()

	sqrt, err = decimalFromStr(mul).ApproxSqrt()

	return sqrt, err
}
