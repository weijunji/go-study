package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMysqlConn(t *testing.T) {
	db := GetMysql()
	mysql, _ := db.DB()
	err := mysql.Ping()
	assert.Nil(t, err)
}
