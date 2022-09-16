package rest

import (
	"github.com/DonggyuLim/uniswap/db"
	"github.com/DonggyuLim/uniswap/pair"
	"github.com/DonggyuLim/uniswap/pool"
	"github.com/DonggyuLim/uniswap/utils"
	"github.com/gin-gonic/gin"
)

type createPairRequest struct {
	AToken pool.Token `json:"AToken"`
	BToken pool.Token `json:"BToken"`
}

func createPair(c *gin.Context) {
	r := createPairRequest{}
	err := c.ShouldBindJSON(&r)

	if err != nil {
		c.String(400, err.Error())
	}
	//이미 페어풀이 있는지 확인
	_, err = db.Get("pair", pair.GetKey(r.AToken.TokenName(), r.BToken.Name))
	if err != nil {
		c.String(400, err.Error())
		return
	}

	p, err := pair.CreatePair(r.AToken, r.BToken)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	db.Add("pair", p.Name, utils.DataToByte(p))
	// pool 을 전체로 보내주는게 맞나? ㅋㅋ
	c.JSON(200, gin.H{
		"message":  "success",
		"pairName": p.GetName(),
		"p":        p,
	})
}
