package main

import (
	"github.com/gin-gonic/gin"
	bptree "github.com/weijunji/bptree/memory"
	"net/http"
)

var trees map[string]*bptree.BPTree = make(map[string]*bptree.BPTree)

func setupRouter() *gin.Engine {
	r := gin.Default()
	tree := r.Group("/bptree")
	{
		bptree.NewBPTree()
		// TODO: get val from tree, return 404 if key doesn't exist
		// GET /bptree/:id?key=123
		// tree.GET("/:id", getFromTree)
		tree.PUT("/:id", putIntoTree)
		// TODO: delete key from tree
		tree.DELETE("/:id", deleteTreeNode)
		// DELETE /bptree/:id  {"key": 123}
		// tree.DELETE("/:id", deleteTreeNode)
	}
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

type TreeKV struct {
	Key uint64 `json:"key" binding:"required"`
	Val string `json:"val" binding:"required"`
}

// insert or update a key into bptree
// create new tree if tree not exists
func putIntoTree(c *gin.Context) {
	id := c.Param("id")
	if _, ok := trees[id]; !ok {
		trees[id] = bptree.NewBPTree()
	}
	tree := trees[id]
	var kv TreeKV
	if err := c.BindJSON(&kv); err != nil {
		c.Status(http.StatusBadRequest)
	}
	if err := tree.Update(kv.Key, kv.Val); err != nil {
		tree.Insert(kv.Key, kv.Val)
	}
	c.Status(http.StatusOK)
}
func deleteTreeNode(c *gin.Context) {
	id := c.Param("id")
	if _, ok := trees[id]; !ok {
		trees[id] = bptree.NewBPTree()
	}
	tree := trees[id]
	var kv TreeKV
	if err := c.BindJSON(&kv); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := tree.Delete(kv.Key); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
