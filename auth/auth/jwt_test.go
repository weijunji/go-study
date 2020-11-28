package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	token, _ := GenerateToken(123, 0, time.Second)
	claims, ok := ParseToken(token)
	assert.True(t, ok)
	assert.Equal(t, claims.ID, uint64(123))
	assert.Equal(t, claims.Role, uint64(0))
	time.Sleep(2 * time.Second)
	_, ok = ParseToken(token)
	assert.False(t, ok)
}
