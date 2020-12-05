package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisCoon(t *testing.T) {
	rdb := GetRedis()
	_, err := rdb.Ping(context.Background()).Result()
	assert.Nil(t, err)
}
