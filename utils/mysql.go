package utils

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var mysql_once sync.Once

// GetDB : return databast instance
func GetMysql() *gorm.DB {
	mysql_once.Do(func() {
		var err error
		if db, err = gorm.Open(mysql.Open(getMysqlSource()), &gorm.Config{}); err != nil {
			panic("Failed to connect mysql")
		}
	})
	return db
}
