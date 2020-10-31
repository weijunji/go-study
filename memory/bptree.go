package memory

import "sync"

const order int = 10

type BPTree struct {
	root     *Node
	nodePool *sync.Pool
}

type Node struct {
	isLeaf  bool
	keys    []uint64
	records []string

	children []*Node
	next     *Node
	prev     *Node
	parent   *Node
}

// func NewBPTree() *BPTree {}
// func (t *BPTree) insert(key uint64, value string) error {}
// func (t *BPTree) delete(key uint64) error {}
// func (t *BPTree) find(key uint64) (string, error) {}
// func (t *BPTree) update(key uint64, value string) error {}
