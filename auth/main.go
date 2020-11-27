package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	. "github.com/weijunji/go-study/auth/utilities"
)

func main() {
	db := GetDB()
	fmt.Printf("db: %+v", db)
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}
