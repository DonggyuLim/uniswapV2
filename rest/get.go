package rest

import (
	"github.com/DonggyuLim/uniswap/db"
	"github.com/DonggyuLim/uniswap/pair"
	"github.com/gin-gonic/gin"
)

type pairRequest struct {
	PairName string `json:"pairName"`
}

func getPair(c *gin.Context) {
	r := pairRequest{}
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	pairByte, err := db.Get("pair", r.PairName)

	if len(pairByte) == 0 || err != nil {
		c.String(400, "Not exsists Pair")
		return
	}
	pair, err := pair.ByteToPair(pairByte)

	if err != nil {
		c.String(400, "Decode error")
		return
	}

	c.JSON(200, gin.H{
		"pair": pair,
	})
}
