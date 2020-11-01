package memory

import "testing"

var testData = []struct {
	op    string
	key   uint64
	value string
}{
	{"insert", 10, "test10"},
	{"insert", 6, "test6"},
	{"insert", 3, "test3"},
	{"insert", 13, "test13"},
	{"insert", 7, "test7"},
	{"insert", 20, "test20"},
	{"insert", 1, "test1"},
}

func TestBPTree(t *testing.T) {
	tree := NewBPTree()
	for _, data := range testData {
		switch data.op {
		case "insert":
			tree.Insert(data.key, data.value)
		}
	}
	// use 'go test . -v' to see the tree
	printTree(tree)
}
