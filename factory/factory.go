package factory

import (
	"log"
	"strings"
)

type address string

// a == b -> true
// a != b -> false
func (a address) isEqual(b address) bool {
	return a == b
}

// a prefix == "0x" -> true
// a prefix != "0x" -> false
func (a address) isAddress() bool {
	return strings.HasPrefix(string(a), "0x")
}

type factory struct {
	pair map[address]address
}

func (f *factory) allPairLength() uint {
	return uint(len(f.pair))
}

func (f *factory) createPair(tokenA, tokenB address) {
	if tokenA.isEqual(tokenB) {
		log.Panic("UniswapV2: IDENTICAL_ADDRESSES'")
	}
	if !tokenA.isAddress() || !tokenB.isAddress() {
		log.Panic("UniswapV2: ZERO_ADDRESS")
	}
	if f.pair[tokenA] == tokenB {
		log.Panicf("%s pair %s is exsits", tokenA, tokenB)
	}
	f.pair[tokenA] = tokenB

}
