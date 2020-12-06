package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	token, _ := GenerateToken(123, 0, time.Second)
	claims, ok := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTIzLCJSb2xlIjowLCJTdGFuZGFyZENsYWltcyI6eyJleHAiOjQ3NjA4NTU2NDN9fQ.UopgmKvymnHxzKHgOKNXNxrDswp-B63mdUwJfsHBrwA")
	fmt.Println(claims)
	assert.True(t, ok)
	assert.Equal(t, claims.ID, uint64(123))
	assert.Equal(t, claims.Role, uint64(0))
	time.Sleep(2 * time.Second)
	_, ok = ParseToken(token)
	assert.False(t, ok)
}
