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
// func (t *BPTree) Insert(key uint64, value string) error {}
// func (t *BPTree) Delete(key uint64) error {}
// func (t *BPTree) Find(key uint64) (string, error) {}
// func (t *BPTree) Update(key uint64, value string) error {}
