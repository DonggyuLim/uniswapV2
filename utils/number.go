package utils

import (
	m "cosmossdk.io/math"
)

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

func ChanegeUintFromInt(a m.Int) m.Uint {
	return m.NewUintFromString(a.String())
}
