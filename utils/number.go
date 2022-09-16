package utils

import "github.com/shopspring/decimal"

// func GetRatio(a, b int64) int64 {
// 	var Ratio int64
// 	if a > b {
// 		Ratio = a / b
// 	} else if a == b {
// 		Ratio = 1
// 	} else {
// 		Ratio = b / a
// 	}
// 	return Ratio
// }

func NewDecimalFromUint(a uint64) decimal.Decimal {
	result := decimal.NewFromInt(int64(a))
	return result
}

func DecimalToUint64(dec decimal.Decimal) uint64 {
	result := dec.BigInt().Uint64()
	return result
}
