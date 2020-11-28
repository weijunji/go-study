package main

import (
	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-study/auth/auth"
)

func main() {
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	anonymousGroup := r.Group("/", auth.AuthMiddleware())
	authGroup := r.Group("/", auth.LoginRequired())
	anonymousGroup.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	auth.SetupRouter(anonymousGroup, authGroup)
	return r
}
