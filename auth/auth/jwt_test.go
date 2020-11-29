package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	token, _ := GenerateToken(123, 0, time.Second)
	claims, ok := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiUm9sZSI6MCwiU3RhbmRhcmRDbGFpbXMiOnsiZXhwIjoxNjA2ODk3MTc5fX0.tXUE5e3-KnIUyx6T9MtKqdzjtZUmn_6ms3dl806Hfko")
	assert.True(t, ok)
	fmt.Println(claims.ID)
	assert.Equal(t, claims.ID, uint64(1))
	assert.Equal(t, claims.Role, uint64(0))
	time.Sleep(2 * time.Second)
	_, ok = ParseToken(token)
	assert.False(t, ok)
}
