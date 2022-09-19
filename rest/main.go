package rest

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const port string = ":8081"

type Document struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Example     string `json:"example,omitempty"`
}

func document(c *gin.Context) {
	data := []Document{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         "/createPair",
			Method:      "POST",
			Description: "Create Pair!",
			Example:     `{name:token-name,symbol:token-symbol,totalSupply:token-supply(uint)}`,
		},
	}
	c.JSON(http.StatusOK, data)
}

func Rest(wg *sync.WaitGroup) {

	r := gin.Default()
	r.GET("/", document)
	r.GET("/pair", getPair)
	r.POST("/createPair", createPair)
	r.POST("/deposit", deposit)
	r.POST("/withdraw", withdraw)
	r.POST("/swap", swap)
	r.Run(port)
	defer wg.Done()
}
