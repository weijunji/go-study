package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Userinfo struct {
	ID   uint64
	Role uint64
}

// AuthMiddleware : only set userinfo if token exist and valid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if claims, ok := ParseToken(token); ok {
			c.Set("userinfo", Userinfo{claims.ID, claims.Role})
		}
	}
}

// LoginRequired : set userinfo and abort the request if not authorized
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if claims, ok := ParseToken(token); ok {
			c.Set("userinfo", Userinfo{claims.ID, claims.Role})
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
