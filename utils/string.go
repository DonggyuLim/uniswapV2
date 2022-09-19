package utils

import (
	"fmt"
	"sort"
)

// usdt,dai => dai:usdt
func GetKey(a, b string) string {
	slice := []string{a, b}
	sort.Strings(slice)
	return fmt.Sprintf("%s:%s", slice[0], slice[1])
}
