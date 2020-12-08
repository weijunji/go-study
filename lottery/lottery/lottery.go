package lottery

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gcache"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastrand"
	"github.com/weijunji/go-study/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RATE_DENOMINATOR = 1000000

// SetupRouter set up lottery router
func SetupRouter(authGroup *gin.RouterGroup) {
	authG := authGroup.Group("/lottery")
	{
		authG.GET("/list", getLotteryList)
		authG.GET("/info/:id", getLotteryInfo)
		authG.GET("/lottery/:id", handleLottery)
	}
}

type AwardInfo struct {
	ID          string `bson:"_id" json:"-"`
	Name        string
	Description string
	Pic         string
	Total       uint32
	DisplayRate uint32 `bson:"displayRate" json:"rate"`
	Rate        uint32 `json:"-"`
	Value       uint32 `json:"-"`
}

type LotteryInfo struct {
	ID          string `bson:"_id"`
	Title       string
	Description string
	Awards      []AwardInfo `bson:"awards"`
	rateSum     uint32
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

var errInvalidID = errors.New("id is not valid objectId")
var errNoSuchLottery = errors.New("no such lottery")

func getInfoByID(id string) (info *LotteryInfo, err error) {
	if infoI, _ := gcache.Get(id); infoI == nil {
		// cache miss
		info = new(LotteryInfo)
		coll := utils.GetMongoDB().Collection("lottery_info")
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			if err := coll.Find(context.Background(), bson.M{"_id": objID}).One(info); err == nil {
				countRate(info)
				gcache.Set(id, info, time.Minute*10)
			} else {
				log.Errorln(err)
				if errors.Is(err, qmgo.ErrNoSuchDocuments) {
					gcache.Set(id, nil, time.Minute*10)
					return nil, errNoSuchLottery
				}
			}
		} else {
			return nil, errInvalidID
		}
	} else {
		// cache hit
		info = infoI.(*LotteryInfo)
		if info == nil {
			return nil, errNoSuchLottery
		}
	}
	return
}

func countRate(info *LotteryInfo) {
	var sum uint32 = 0
	for _, award := range info.Awards {
		sum += award.Rate
	}
	info.rateSum = sum
}

func getLotteryInfo(c *gin.Context) {
	id := c.Param("id")
	if info, err := getInfoByID(id); err != nil {
		if errors.Is(err, errInvalidID) {
			c.Status(http.StatusBadRequest)
		} else if errors.Is(err, errNoSuchLottery) {
			c.Status(http.StatusNotFound)
		} else {
			log.Errorln(err)
			c.Status(http.StatusInternalServerError)
		}
	} else {
		c.JSON(http.StatusOK, info)
	}
}

func handleLottery(c *gin.Context) {
	id := c.Param("id")
	if info, err := getInfoByID(id); err != nil {
		if errors.Is(err, errInvalidID) {
			c.Status(http.StatusBadRequest)
		} else if errors.Is(err, errNoSuchLottery) {
			c.Status(http.StatusNotFound)
		} else {
			log.Errorln(err)
			c.Status(http.StatusInternalServerError)
		}
	} else {
		result := processLottery(info)
		if result != nil {
			if success := deleteOne(info.ID, result); success {
				c.JSON(http.StatusOK, gin.H{"code": 1, "message": "恭喜中奖", "award": result})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "未中奖"})
	}
}

func deleteLotteryTime(id uint64) {
	// read redis -> hit -> try delete
	//     |-> not hit -> read mysql -> found -> set redis -> try delete
	//                        |-> not found -> create -|
	// 使用redis过期事件来监听 -> mq -> 写回
}

func processLottery(info *LotteryInfo) *AwardInfo {
	randNum := fastrand.Uint32n(RATE_DENOMINATOR)
	if randNum < info.rateSum {
		var count uint32 = 0
		for _, award := range info.Awards {
			count += award.Rate
			if count >= randNum {
				return &award
			}
		}
	}
	return nil
}

func deleteOne(lotteryID string, award *AwardInfo) bool {
	if award.Value > 100 {
		// high value in mysql
		db, _ := utils.GetMysql().DB()
		if tx, err := db.Begin(); err == nil {
			var remain uint32
			tx.QueryRow("SELECT remain FROM awards WHERE award = ?", award.ID).Scan(&remain)
			logrus.WithFields(logrus.Fields{"award": award.ID, "remain": remain}).Debugf("Try delete one from mysql")
			if remain > 0 {
				if _, err := tx.Exec("UPDATE awards SET remain = ? WHERE award = ?", remain-1, award.ID); err == nil {
					tx.Commit()
					return true
				}
			}
			tx.Rollback()
		}
	} else {
		// low value in redis
		db := utils.GetRedis()
		key := "awards:" + lotteryID + ":" + award.ID
		if remain, err := db.IncrBy(context.Background(), key, -1).Result(); err == nil {
			return remain >= 0
		}
	}
	return false
}
