package rest

import (
	"sync"

	"github.com/gin-gonic/gin"
)

const port string = ":7000"

func Rest(wg *sync.WaitGroup) {

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "hello world")
	})
	r.POST("/createPair", createPair)
	r.Run(port)
	defer wg.Done()
}
