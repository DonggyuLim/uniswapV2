package rest

import (
	"errors"

	"github.com/DonggyuLim/uniswap/db"
	"github.com/DonggyuLim/uniswap/pair"
	"github.com/DonggyuLim/uniswap/utils"
)

func savePair(p pair.Pair) {
	db.Add("pair", p.GetName(), utils.DataToByte(p))
}
func loadPair(pairName string) (pair.Pair, error) {
	pairByte, err := db.Get("pair", pairName)
	if len(pairByte) == 0 || err != nil {
		return pair.Pair{}, errors.New("not exsits Pair")
	}
	p, err := pair.ByteToPair(pairByte)
	if err != nil {
		return pair.Pair{}, errors.New("decode fail")
	}
	return p, nil
}
