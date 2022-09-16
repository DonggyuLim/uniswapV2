package pool

// import (
// 	"testing"

// 	"github.com/shopspring/decimal"
// )

// func TestDeposit(t *testing.T) {
// 	type test struct {
// 		px int64
// 		py int64
// 		pc float64
// 		a  int64
// 		b  int64
// 		lp float64
// 	}
// 	tests := []test{
// 		{
// 			px: 200,
// 			py: 400,
// 			pc: 282,
// 			a:  10,
// 			b:  20,
// 			lp: 14.142135623730951,
// 		},
// 		{
// 			px: 3000,
// 			py: 2000,
// 			pc: 2449.48974278,
// 			a:  100,
// 			b:  300,
// 			lp: 122.47448713915891,
// 		},
// 		{
// 			px: 200,
// 			py: 300,
// 			pc: 244.948974278,
// 			a:  10,
// 			b:  30,
// 			lp: 12.24744871391589,
// 		},
// 		{
// 			px: 300,
// 			py: 200,
// 			pc: 244.948974278,
// 			a:  1000,
// 			b:  300,
// 			lp: 367.4234614174767,
// 		},
// 		{
// 			px: 300,
// 			py: 1000,
// 			pc: 547.722557505,
// 			a:  1000,
// 			b:  2000,
// 			lp: 1095.4451150103323,
// 		},
// 		{
// 			px: 300,
// 			py: 1000,
// 			pc: 547.722557505,
// 			a:  2000,
// 			b:  1000,
// 			lp: 547.7225575051662,
// 		},
// 	}
// 	t.Run("Deposit", func(t *testing.T) {
// 		for _, tc := range tests {
// 			tokenA := Token{
// 				Name:    "bitcoin",
// 				Balance: newInt(tc.px),
// 			}

// 			tokenB := Token{
// 				Name:    "Atom",
// 				Balance: newInt(tc.py),
// 			}
// 			pc := Token{
// 				Name:    "uni-ba",
// 				Balance: decimal.NewFromFloat(tc.pc),
// 			}
// 			pool := CreatePool(tokenA, tokenB, pc)
// 			tokenA.Balance = newInt(tc.a)
// 			tokenB.Balance = newInt(tc.b)

// 			pc = pool.Deposit(tokenA, tokenB)
// 			if pc.Balance.Cmp(decimal.NewFromFloat(tc.lp)) != 0 {
// 				t.Errorf("Expected %v , got %v", tc.lp, pc.Balance.Abs())
// 			}
// 		}

// 	})

// }

// func TestWithDraw(t *testing.T) {
// 	type test struct {
// 		px int64
// 		py int64
// 		pc float64
// 		lp int64
// 		x  float64
// 		y  float64
// 	}

// 	tests := []test{
// 		{
// 			px: 10000,
// 			py: 20000,
// 			pc: 14142.1356237,
// 			lp: 1000,
// 			x:  707.106781188095,
// 			y:  1414.21356237619,
// 		},
// 		{
// 			px: 11113,
// 			py: 20302,
// 			pc: 15020.5234929,
// 			lp: 350,
// 			x:  258.9490307603820541,
// 			y:  473.0660687930600614,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run("TestWithDraw!", func(t *testing.T) {
// 			tokenA := Token{
// 				Name:    "bitcoin",
// 				Balance: newInt(tc.px),
// 			}

// 			tokenB := Token{
// 				Name:    "Atom",
// 				Balance: newInt(tc.py),
// 			}
// 			pc := Token{
// 				Name:    "uni-ba",
// 				Balance: newFloat(tc.pc),
// 			}
// 			pool := CreatePool(tokenA, tokenB, pc)
// 			lp := Token{
// 				Name:    "uni-ba",
// 				Balance: newInt(tc.lp),
// 			}
// 			x, y := pool.WithDraw(lp)

// 			if newFloat(tc.x).Cmp(x.Balance) != 0 {
// 				t.Errorf("Expected %v got %v", tc.x, x.Balance)
// 			}
// 			if newFloat(tc.y).Cmp(y.Balance) != 0 {
// 				t.Errorf("Expected %v got %v", tc.y, y.Balance)
// 			}
// 		})
// 	}
// }
