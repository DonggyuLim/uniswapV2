package utils

func GetRatio(a, b int64) int64 {
	var Ratio int64
	if a > b {
		Ratio = a / b
	} else if a == b {
		Ratio = 1
	} else {
		Ratio = b / a
	}
	return Ratio
}
