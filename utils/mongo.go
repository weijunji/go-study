package utils

import (
	"context"
	"sync"

	"github.com/qiniu/qmgo"
)

var mongo *qmgo.Database
var mongoOnce sync.Once

// GetMongoDB return a qmgo client
func GetMongoDB() *qmgo.Database {
	mongoOnce.Do(func() {
		cli, err := qmgo.Open(context.Background(), getMongoConfig())
		if err != nil {
			panic("Failed to connect to MongoDB")
		}
		mongo = cli.Database
	})
	return mongo
}
