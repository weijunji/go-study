package lottery

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
	log "github.com/sirupsen/logrus"
	"github.com/weijunji/go-study/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetupRouter set up lottery router
func SetupRouter(authGroup *gin.RouterGroup) {
	authG := authGroup.Group("/lottery")
	{
		authG.GET("/list", getLotteryList)
		authG.GET("/info/:id", getLotteryInfo)
		authG.GET("/lottery/:id", processLottery)
	}
}

type AwardInfo struct {
	Name        string
	Description string
	Pic         string
	Total       int
	DisplayRate int `bson:"displayRate"`
	Rate        int
	Value       int
}

type LotteryInfo struct {
	ID          string `bson:"_id"`
	Title       string
	Description string
	Awards      []AwardInfo `bson:"awards"`
}

func getLotteryList(c *gin.Context) {
	type SimpleLotteryInfo struct {
		ID          string `bson:"_id"`
		Title       string
		Description string
	}
	coll := utils.GetMongoDB().Collection("lottery_info")
	infos := []SimpleLotteryInfo{}
	if err := coll.Find(context.Background(), bson.M{}).All(&infos); err != nil {
		log.Errorln(err)
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, &infos)
	}
}

func getLotteryInfo(c *gin.Context) {
	id := c.Param("id")
	coll := utils.GetMongoDB().Collection("lottery_info")
	var info LotteryInfo
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		if err := coll.Find(context.Background(), bson.M{"_id": objID}).One(&info); err != nil {
			log.Errorln(err)
			if errors.Is(err, qmgo.ErrNoSuchDocuments) {
				c.Status(http.StatusNotFound)
			} else {
				c.Status(http.StatusInternalServerError)
			}
		} else {
			c.JSON(http.StatusOK, &info)
		}
	} else {
		log.Errorln(err)
		c.Status(http.StatusBadRequest)
	}
}

func processLottery(c *gin.Context) {
	// id := c.Param("id")
}
