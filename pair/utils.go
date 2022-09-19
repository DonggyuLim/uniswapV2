package pair

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/DonggyuLim/uniswap/pool"
)

func checkSameToken(tokenA, tokenB pool.Token) error {
	if tokenA.Name == tokenB.Name {
		return errors.New("same tokens can't become pair")
	}
	return nil
}

// usdt,dai => dai:usdt
func MakeKey(a, b string) string {
	slice := []string{a, b}
	sort.Strings(slice)
	return fmt.Sprintf("%s:%s", slice[0], slice[1])
}

// usdt,dai => dnuLP
func MakeSymbol(a, b string) string {
	aslice := strings.Split(a, "")
	bslice := strings.Split(b, "")
	result := fmt.Sprintf("%sN%sLP", aslice[0], bslice[0])
	return result
}

func ByteToPair(data []byte) (Pair, error) {
	var pair Pair
	decoder := gob.NewDecoder(bytes.NewBuffer((data)))
	err := decoder.Decode(&pair)
	if err != nil {
		return pair, err
	}
	return pair, nil
}
