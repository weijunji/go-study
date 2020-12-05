package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoConn(t *testing.T) {
	GetMongoDB()
	coll := mongo.Collection("myuser")
	count, _ := coll.Find(context.Background(), bson.M{"name": "test"}).Count()
	assert.NotEqual(t, 0, count)
}
