package main

import (
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-study/auth/auth"
	"github.com/weijunji/go-study/lottery/lottery"
)

func main() {
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	authGroup := r.Group("/", auth.LoginRequired())
	lottery.SetupRouter(authGroup)
	return r
}
