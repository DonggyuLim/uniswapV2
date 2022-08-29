package pool

import (
	"math/big"
	"testing"
)

func TestReserve(t *testing.T) {
	tokenA := token{
		Name:    "bitcoin",
		Balance: big.NewInt(10000),
	}

	tokenB := token{
		Name:    "Atom",
		Balance: big.NewInt(200000),
	}
	pc := token{
		Name:    "uni-ba",
		Balance: big.NewInt(44721),
	}
	pool := CreatePool(tokenA, tokenB, pc)
	t.Run("Reserve offer reserve balance", func(t *testing.T) {
		x, y := pool.Reserve()
		if x != tokenA.Balance {
			t.Errorf("Expected %s, got %s", tokenA.Balance, x)
		}
		if y != tokenB.Balance {
			t.Errorf("Expected %s,got %s", tokenB.Balance, y)
		}
	})

}

func TestK(t *testing.T) {
	type test struct {
		x int64
		y int64
		k int64
	}
	tests := []test{
		{
			x: 10,
			y: 1000,
			k: 10000,
		},
		{
			x: 253,
			y: 3000,
			k: 759000,
		},
		{
			x: 100,
			y: 2310,
			k: 231000,
		},
		{
			x: 12301,
			y: 401113,
			k: 4934091013,
		},
	}
	for _, tc := range tests {

		tokenA := token{
			Name:    "bitcoin",
			Balance: big.NewInt(tc.x),
		}

		tokenB := token{
			Name:    "Atom",
			Balance: big.NewInt(tc.y),
		}
		ps := token{
			Name:    "uni-BA",
			Balance: big.NewInt(tc.k),
		}
		pool := CreatePool(tokenA, tokenB, ps)
		got := pool.K()

		if big.NewInt(tc.k).Cmp(got) != 0 {
			t.Errorf("Expected %v and got %v", tc.k, got)
		}
	}
}