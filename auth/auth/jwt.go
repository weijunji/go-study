package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const secretKey string = "FlFhhEDskm68mOAbi!yTryI0KleXlgJ@" // replace this before deployment

type LoginClaims struct {
	ID             uint64
	Role           uint64
	StandardClaims jwt.StandardClaims
}

// Valid : implement jwt.Claims, check expire time
func (claims LoginClaims) Valid() error {
	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return errors.New("Token already expired")
	}
	return nil
}

// GenerateToken : generate token with 3 days
func GenerateToken(uid uint64, role uint64, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LoginClaims{
		ID:   uid,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	})
	return token.SignedString([]byte(secretKey))
}

// ParseToken : parse token return claims
func ParseToken(tokenStr string) (claims *LoginClaims, ok bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok = token.Claims.(*LoginClaims); ok && token.Valid {
		return
	}
	return nil, false
}
