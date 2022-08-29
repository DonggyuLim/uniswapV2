package pool

import (
	"math/big"
	"testing"
)

func TestDeposit(t *testing.T) {
	type test struct {
		px int64
		py int64
		pc int64
		a  int64
		b  int64
		lp int64
	}
	tests := []test{
		{
			px: 200,
			py: 400,
			pc: 282,
			a:  10,
			b:  20,
			lp: 14,
		},
		{
			px: 2000,
			py: 4000,
			pc: 2828,
			a:  30,
			b:  70,
			lp: 45,
		},
		{
			px: 14123,
			py: 1512312,
			pc: 146145,
			a:  141,
			b:  15098,
			lp: 1459,
		},
		// {
		// 	px: 1011,
		// 	py: 1203, //1.18 일때를 해결해야함 전부다 리팩토링?
		// 	pc: 1102,
		// 	a:  123,
		// 	b:  10405,
		// 	lp: 133,
		// },
	}
	t.Run("Deposit", func(t *testing.T) {
		for _, tc := range tests {
			tokenA := token{
				Name:    "bitcoin",
				Balance: big.NewInt(tc.px),
			}

			tokenB := token{
				Name:    "Atom",
				Balance: big.NewInt(tc.py),
			}
			pc := token{
				Name:    "uni-ba",
				Balance: big.NewInt(tc.pc),
			}
			pool := CreatePool(tokenA, tokenB, pc)
			tokenA.Balance = big.NewInt(tc.a)
			tokenB.Balance = big.NewInt(tc.b)

			pc.Balance = pool.Deposit(tokenA, tokenB).Balance
			if pc.Balance.Cmp(big.NewInt(tc.lp)) != 0 {
				t.Errorf("Expected %v , got %s", tc.lp, pc.Balance)
			}
		}

	})

}

func TestWithDraw(t *testing.T) {
	type test struct {
		px int64
		py int64
		pc int64
		lp int64
		x  int64
		y  int64
	}

	tests := []test{
		{
			px: 10000,
			py: 20000,
			pc: 14142,
			lp: 1000,
			x:  700,
			y:  1400,
		},
		{
			px: 11113,
			py: 20302,
			pc: 15020,
			lp: 351,
			x:  222,
			y:  406,
		},
	}

	for _, tc := range tests {
		t.Run("TestWithDraw!", func(t *testing.T) {
			tokenA := token{
				Name:    "bitcoin",
				Balance: big.NewInt(tc.px),
			}

			tokenB := token{
				Name:    "Atom",
				Balance: big.NewInt(tc.py),
			}
			pc := token{
				Name:    "uni-ba",
				Balance: big.NewInt(tc.pc),
			}
			pool := CreatePool(tokenA, tokenB, pc)
			lp := token{
				Name:    "uni-ba",
				Balance: big.NewInt(tc.lp),
			}
			x, y := pool.WithDraw(lp)

			if big.NewInt(tc.x).Cmp(x.Balance) != 0 {
				t.Errorf("Expected %v got %v", tc.x, x.Balance)
			}
			if big.NewInt(tc.y).Cmp(y.Balance) != 0 {
				t.Errorf("Expected %v got %v", tc.y, y.Balance)
			}
		})
	}
}
