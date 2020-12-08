package lottery

import "github.com/gin-gonic/gin"

// SetupRouter set up lottery router
func SetupRouter(authGroup *gin.RouterGroup) {
	authG := authGroup.Group("/lottery")
	{
		authG.GET("/list", getLotteryList)
		authG.GET("/{id}", getLotteryInfo)
		authG.GET("/lottery/{id}", processLottery)
	}
}

func getLotteryList(c *gin.Context) {

}

func getLotteryInfo(c *gin.Context) {

}

func processLottery(c *gin.Context) {

}
