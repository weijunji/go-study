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
	type UserRegistration struct {
		Username string `json:"username" binding:"required"`
		Passwd1  string `json:"password1" binding:"required"`
		Passwd2  string `json:"password2" binding:"required"`
	}
	userRegister := UserRegistration{}
	if err := c.BindJSON(&userRegister); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if userRegister.Passwd1 != userRegister.Passwd2 {
		c.Status(http.StatusBadRequest)
		return
	}
	//the username has been existed
	tx := util.GetDB().Where("username=?", userRegister.Username).First(&User{})
	if tx.RowsAffected != 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	userInfo := User{
		Username: userRegister.Username,
		Password: encryptPassword(userRegister.Passwd1),
		Role:     0,
	}
	tx = util.GetDB().Create(&userInfo)
	if err := tx.Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//regist success,return token
	token, err := GenerateToken(userInfo.ID, userInfo.Role, 3*24*time.Hour)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}

func updateProfile(c *gin.Context) {
	// val, _ := c.Get("userinfo")
	// userinfo, _ := val.(Userinfo)
	// TODO: update user profile
	val, _ := c.Get("userinfo")
	userinfo := val.(Userinfo)
	type UpdateBody struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		AvatarURL string `json:"avatarURL"`
		Introduce string `json:"introduce"`
		Email     string `json:"email"`
	}
	updateBody := UpdateBody{}
	if err := c.BindJSON(&updateBody); err != nil {
		c.Status(http.StatusBadRequest)
	}
	if updateBody.Password != "" {
		updateBody.Password = encryptPassword(updateBody.Password)
	}
	user := User{
		ID:        userinfo.ID,
		Username:  updateBody.Username,
		Password:  updateBody.Password,
		AvatarURL: updateBody.AvatarURL,
		Introduce: updateBody.Introduce,
		Email:     updateBody.Email,
	}
	if err := util.GetDB().Model(&user).Updates(user).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, struct {
		ID        uint64 `json:"id"`
		Username  string `json:"username"`
		Role      uint64 `json:"role"`
		AvatarURL string `json:"avatarURL"`
		Introduce string `json:"introduce"`
		Email     string `json:"email"`
	}{
		user.ID,
		user.Username,
		user.Role,
		user.AvatarURL,
		user.Introduce,
		user.Email,
	})
}

func getProfile(c *gin.Context) {
	// val, _ := c.Get("userinfo")
	// userinfo, _ := val.(Userinfo)
	// TODO: get user profile
	val, _ := c.Get("userinfo")
	userinfo, _ := val.(Userinfo)
	user := User{ID: userinfo.ID}
	if err := util.GetDB().First(&user).Error; err != nil {
		c.Status(http.StatusOK)
		return
	}
	c.JSON(http.StatusOK, struct {
		ID        uint64 `json:"id"`
		Username  string `json:"username"`
		Role      uint64 `json:"role"`
		AvatarURL string `json:"avatarURL"`
		Introduce string `json:"introduce"`
		Email     string `json:"email"`
	}{
		user.ID,
		user.Username,
		user.Role,
		user.AvatarURL,
		user.Introduce,
		user.Email,
	})
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
