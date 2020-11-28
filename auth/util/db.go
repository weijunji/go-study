package util

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

// GetDB : return databast instance
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		if db, err = gorm.Open(sqlite.Open("auth.sqlite"), &gorm.Config{}); err != nil {
			panic("failed to connect database")
		}
	})
	return db
}
