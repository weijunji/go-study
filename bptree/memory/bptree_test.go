package memory

import (
	"fmt"
	"testing"
)

var testData = []struct {
	op    string
	key   uint64
	value string
}{
	{"update", 100, "test100"},
	{"insert", 10, "test10"},
	{"insert", 6, "test6"},
	{"insert", 3, "test3"},
	{"insert", 13, "test13"},
	{"insert", 7, "test7"},
	{"insert", 20, "test20"},
	{"insert", 1, "test1"},
	// kdjlyy
	{"update", 3, "test3_updated"},
	{"update", 16, "test16_updated"},
	{"update", 6, "test6_updated"},
	{"find", 7, ""},
	{"find", 100, ""},
}

func TestBPTree(t *testing.T) {
	tree := NewBPTree()
	for _, data := range testData {
		switch data.op {
		case "insert":
			tree.Insert(data.key, data.value)
		case "find":
			{ // kdjlyy
				record, err := tree.Find(data.key)
				if err != nil {
					fmt.Printf("%s in Find(%d)!!!\n", err, data.key)
				} else {
					fmt.Printf("Find(%d) = %s\n", data.key, record)
				}
			}
		case "update":
			{ // kdjlyy
				err := tree.Update(data.key, data.value)
				if err != nil {
					fmt.Printf("%s in Update(%d, %s)!!!\n", err, data.key, data.value)
				}
			}

		}
	}

	PrintTree(tree)

	tree.Delete(13)
	PrintTree(tree)
	// use 'go test . -v' to see the tree
}
