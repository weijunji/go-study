package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bptree "github.com/weijunji/bptree/memory"
)

var trees map[string]*bptree.BPTree = make(map[string]*bptree.BPTree)

func setupRouter() *gin.Engine {
	r := gin.Default()
	tree := r.Group("/bptree")
	{
		bptree.NewBPTree()
		// TODO: get val from tree, return 404 if key doesn't exist
		// GET /bptree/:id?key=123
		tree.GET("/:id", getFromTree)
		tree.PUT("/:id", putIntoTree)
		tree.DELETE("/:id", deleteTreeNode)
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
		return
	}
	if err := tree.Update(kv.Key, kv.Val); err != nil {
		tree.Insert(kv.Key, kv.Val)
	}
	c.Status(http.StatusOK)
}

type TreeK struct {
	Key uint64 `json:"key" binding:"required"`
}

func getFromTree(c *gin.Context) {
	id := c.Param("id")
	var tree *bptree.BPTree
	var ok bool
	if tree, ok = trees[id]; !ok {
		c.Status(http.StatusNotFound)
		return
	}
	var key TreeK
	if err := c.BindJSON(&key); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if val, err := tree.Find(key.Key); err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"value": val,
		})
	}
}

func deleteTreeNode(c *gin.Context) {
	id := c.Param("id")
	var tree *bptree.BPTree
	var ok bool
	if tree, ok = trees[id]; !ok {
		c.Status(http.StatusNotFound)
		return
	}
	var key TreeK
	if err := c.BindJSON(&key); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := tree.Delete(key.Key); err != nil {
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
