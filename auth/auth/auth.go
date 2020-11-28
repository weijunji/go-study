package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/weijunji/go-study/auth/util"
	"gorm.io/gorm"
)

// User : struct for user
type User struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"type:varchar(16);uniqueIndex;not null"`
	Password  string `gorm:"type:char(32);not null" json:"-"`
	Role      uint64 // 0: admin, 1: normal
	AvatarURL string
	Introduce string
	Email     string
}

const salt string = "3Wc6cX6A"

// encrypt the password for store
// IMPORTANT!! PASSWORD SHOULD BE ENCRYPTED BEFORE STORE INTO THE DATABASE
func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

// SetupRouter : set up auth router
func SetupRouter(anonymousGroup *gin.RouterGroup, authGroup *gin.RouterGroup) {
	util.GetDB().AutoMigrate(&User{})
	anonyAuthG := anonymousGroup.Group("/auth")
	{
		anonyAuthG.POST("/login", login)
		anonyAuthG.POST("/register", register)
	}
	authG := authGroup.Group("/auth")
	{
		authG.PUT("/profile", updateProfile)
		authG.GET("/profile", getProfile)
	}
}

func register(c *gin.Context) {
	// TODO: register
}

func updateProfile(c *gin.Context) {
	// val, _ := c.Get("userinfo")
	// userinfo, _ := val.(Userinfo)
	// TODO: update user profile
}

func getProfile(c *gin.Context) {
	// val, _ := c.Get("userinfo")
	// userinfo, _ := val.(Userinfo)
	// TODO: get user profile
}

// Login : login via username and password
func login(c *gin.Context) {
	type LoginBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var body LoginBody
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// get user from db
	db := util.GetDB()
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusUnauthorized)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	// check password
	if encryptPassword(body.Password) != user.Password {
		c.Status(http.StatusUnauthorized)
		return
	}

	// generate token
	token, err := GenerateToken(user.ID, user.Role, 3*24*time.Hour)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(200, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
