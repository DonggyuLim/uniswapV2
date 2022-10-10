package rest

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"github.com/DonggyuLim/uniswap/client"
	"github.com/DonggyuLim/uniswap/db"
	"github.com/DonggyuLim/uniswap/pair"
	"github.com/DonggyuLim/uniswap/pool"
	"github.com/DonggyuLim/uniswap/utils"
	"github.com/gin-gonic/gin"
)

type AccountResponse struct {
	TokenName string `json:"tokenName"`
	Account   string `json:"account"`
	Amount    string `json:"amount"`
}

type createPairRequest struct {
	XToken pool.Token `json:"XToken"`
	YToken pool.Token `json:"YToken"`
}

func AccountisEqaul(a, b pool.Token) error {
	if a.Account != b.Account {
		return errors.New("x,y account must Equal")
	}
	return nil
}

func createPair(c *gin.Context) {
	r := createPairRequest{}
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	err = AccountisEqaul(r.XToken, r.YToken)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	address := r.XToken.Account
	//이미 페어풀이 있는지 확인
	pByte, _ := db.Get("pair", utils.GetKey(r.XToken.GetTokenName(), r.YToken.GetTokenName()))

	if len(pByte) != 0 || err != nil {
		pair, _ := pair.ByteToPair(pByte)
		fmt.Println(pair)
		c.String(400, fmt.Sprintf("Exists Pair, pair is %v", pair.GetName()))
		return
	}

	p, err := pair.CreatePair(r.XToken, r.YToken)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	_, err = client.GetClient().Approve(r.XToken.GetTokenName(), r.XToken.Account, "0xuni", r.XToken.GetBalance().String())
	if err != nil {
		c.String(400, err.Error())
		return
	}

	_, err = client.GetClient().Approve(r.YToken.GetTokenName(), r.YToken.Account, "0xuni", r.YToken.GetBalance().String())
	if err != nil {
		c.String(400, err.Error())
		return
	}

	err = p.Pool.Deposit(address, r.XToken, r.YToken)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	// pool 을 전체로 보내주는게 맞나? ㅋㅋ
	c.JSON(200, gin.H{
		"message": "success",
		"p":       p,
	})
}

type depositRequest struct {
	PairName string     `json:"pairName"`
	XToken   pool.Token `json:"XToken"`
	YToken   pool.Token `json:"YToken"`
}

func deposit(c *gin.Context) {
	r := depositRequest{}
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	err = AccountisEqaul(r.XToken, r.YToken)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	address := r.XToken.Account
	p, err := loadPair(r.PairName)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	err = p.Pool.Deposit(address, r.XToken, r.YToken)
	fmt.Println(err)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	savePair(p)
	c.JSON(200, gin.H{
		"message": "success",
		"Balance": p.Pool.LP.BalanceOf(address),
	})
}

type withDrawRequest struct {
	PairName string `json:"pairname"`
	Account  string `json:"account"`
	Amount   string `json:"amount"`
}

func withdraw(c *gin.Context) {
	r := withDrawRequest{}
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	p, err := loadPair(r.PairName)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	err = p.Pool.WithDraw(r.Account, math.NewUintFromString(r.Account))
	if err != nil {
		c.String(400, fmt.Sprintf("Fail Withdraw why? %s", err.Error()))
	}
	savePair(p)
	c.JSON(200, gin.H{
		"meessage": "success",
		"Account": AccountResponse{
			TokenName: p.Pool.GetLPname(),
			Account:   r.Account,
			Amount:    p.Pool.LP.BalanceOf(r.Account).String(),
		},
	})

}

type swapRequest struct {
	PairName  string `json:"pairName"`
	TokenName string `json:"tokenName"`
	Account   string `json:"account"`
	Amount    string `json:"amount"`
}

func swap(c *gin.Context) {
	r := swapRequest{}
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	p, err := loadPair(r.PairName)
	if err != nil {
		c.String(400, "Pair name isn't exsits")
	}
	err = p.Pool.Swap(r.TokenName, r.Account, math.NewUintFromString(r.Amount))
	if err != nil {
		c.String(400, err.Error())
	}
	savePair(p)
	var tokenName string
	if r.TokenName == p.Pool.X.Name {
		tokenName = p.Pool.GetYName()
	} else {
		tokenName = p.Pool.GetXName()
	}
	amount, err := client.GetClient().GetBalance(tokenName, r.Account)
	if err != nil {
		c.String(400, err.Error())
	}
	c.JSON(200, gin.H{
		"account": AccountResponse{
			TokenName: tokenName,
			Account:   r.Account,
			Amount:    amount,
		},
	})
}
