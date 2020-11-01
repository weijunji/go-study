package memory

import (
	"errors"
	"fmt"
	"sort"
)

const order int = 6
const half int = (order + 1) / 2

// ErrKeyExists : key already exists
var ErrKeyExists = errors.New("Key exists")

// ErrKeyNotFound : key not found
var ErrKeyNotFound = errors.New("Key not found")

// BPTree : bptree
type BPTree struct {
	root      *Node
	firstLeaf *Node
}

// Node : node of bptree
type Node struct {
	isLeaf  bool
	keys    []uint64
	records []string

	children []*Node
	next     *Node
	prev     *Node
	parent   *Node
}

func getIndex(keys []uint64, key uint64) int {
	idx := sort.Search(len(keys), func(i int) bool {
		return key <= keys[i]
	})
	return idx
}

// print tree for debug
func printTree(tree *BPTree) {
	queue := make([]*Node, 0)
	if tree.root != nil {
		queue = append(queue, tree.root)
	}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.isLeaf {
			fmt.Printf("\n--LEAF--%p--------------------------\n", node)
			for idx, key := range node.keys {
				fmt.Printf("%d\t%s\n", key, node.records[idx])
			}
		} else {
			fmt.Printf("\n--NODE--%p--------------------------\n", node)
			for idx, key := range node.keys {
				fmt.Printf("%d\t%p\n", key, node.children[idx])
				queue = append(queue, node.children[idx])
			}
		}
	}
}

// NewBPTree : return a empty bptree
func NewBPTree() *BPTree {
	return &BPTree{}
}

func (t *BPTree) findLeaf(key uint64) *Node {
	p := t.root

	for !p.isLeaf {
		/*
			idx := sort.Search(len(p.keys), func(i int) bool {
				return key <= p.keys[i]
			})
		*/
		idx := getIndex(p.keys, key)
		if idx == len(p.keys) {
			idx = len(p.keys) - 1
		}
		p = p.children[idx]
	}
	return p
}

func insertKeyValIntoLeaf(n *Node, key uint64, rec string) (int, error) {
	idx := getIndex(n.keys, key)
	/*
		idx := sort.Search(len(n.keys), func(i int) bool {
			return key <= n.keys[i]
		})
	*/
	if idx < len(n.keys) && n.keys[idx] == key {
		return 0, ErrKeyExists
	}

	n.keys = append(n.keys, key)
	n.records = append(n.records, rec)
	for i := len(n.keys) - 1; i > idx; i-- {
		n.keys[i] = n.keys[i-1]
		n.records[i] = n.records[i-1]
	}
	n.keys[idx] = key
	n.records[idx] = rec
	return idx, nil
}

func (t *BPTree) updateLastParentKey(leaf *Node) {
	key := leaf.keys[len(leaf.keys)-1]
	updateNode := leaf.parent

	p := leaf
	idx := len(leaf.keys) - 1

	for updateNode != nil && idx == len(p.keys)-1 {
		for i, v := range updateNode.children {
			if v == p {
				idx = i
				break
			}
		}
		updateNode.keys[idx] = key
		updateNode = updateNode.parent
		p = updateNode
	}
}

func (t *BPTree) splitLeafIntoTowleaves(leaf *Node, newLeaf *Node) {
	for i := half; i <= order; i++ {
		newLeaf.keys = append(newLeaf.keys, leaf.keys[i])
		newLeaf.records = append(newLeaf.records, leaf.records[i])
	}

	// adjust relation
	leaf.keys = leaf.keys[:half]
	leaf.records = leaf.records[:half]

	newLeaf.next = leaf.next
	leaf.next = newLeaf
	newLeaf.prev = leaf

	newLeaf.parent = leaf.parent

	if newLeaf.next != nil {
		newLeaf.next.prev = newLeaf
	}
}

func insertIntoNode(parent *Node, idx int, left *Node, key uint64, right *Node) {
	parent.keys = append(parent.keys, key)
	for i := len(parent.keys) - 1; i > idx; i-- {
		parent.keys[i] = parent.keys[i-1]
	}
	parent.keys[idx] = key

	if idx == len(parent.children) {
		parent.children = append(parent.children, right)
		return
	}
	tmpChildren := append([]*Node{}, parent.children[idx+1:]...)
	parent.children = append(append(parent.children[:idx+1], right), tmpChildren...)
}

func (t *BPTree) insertIntoNodeAfterSplitting(oldNode *Node) {
	newNode := &Node{}
	newNode.isLeaf = false

	for i := half; i <= order; i++ {
		newNode.children = append(newNode.children, oldNode.children[i])
		newNode.keys = append(newNode.keys, oldNode.keys[i])
		// update new_node children relation
		child := oldNode.children[i]
		child.parent = newNode
	}
	newNode.parent = oldNode.parent

	oldNode.children = oldNode.children[:half]
	oldNode.keys = oldNode.keys[:half]

	newNode.next = oldNode.next
	oldNode.next = newNode
	newNode.prev = oldNode

	if newNode.next != nil {
		newNode.next.prev = newNode
	}

	t.insertIntoParent(oldNode.parent, oldNode, oldNode.keys[len(oldNode.keys)-1], newNode)
}

func (t *BPTree) insertIntoParent(parent *Node, left *Node, key uint64, right *Node) {
	if parent == nil {
		root := &Node{}
		root.isLeaf = false
		root.keys = append(root.keys, left.keys[len(left.keys)-1])
		root.keys = append(root.keys, right.keys[len(right.keys)-1])
		root.children = append(root.children, left)
		root.children = append(root.children, right)
		left.parent = root
		right.parent = root

		t.root = root
		return
	}

	idx := getIndex(parent.keys, key)
	insertIntoNode(parent, idx, left, key, right)

	if len(parent.keys) > order {
		t.insertIntoNodeAfterSplitting(parent)
	}
}

func (t *BPTree) insertIntoLeaf(key uint64, value string) error {
	var (
		leaf *Node
		err  error
		idx  int
	)
	leaf = t.findLeaf(key)

	if idx, err = insertKeyValIntoLeaf(leaf, key, value); err != nil {
		return err
	}

	if idx == len(leaf.keys)-1 && leaf.parent != nil {
		t.updateLastParentKey(leaf)
	}

	// insert finish
	if len(leaf.keys) <= order {
		return nil
	}

	// split leaf so new leaf node
	newLeaf := &Node{}
	newLeaf.isLeaf = true
	t.splitLeafIntoTowleaves(leaf, newLeaf)

	// insert split key into parent
	t.insertIntoParent(leaf.parent, leaf, leaf.keys[len(leaf.keys)-1], newLeaf)
	return nil
}

// Insert : insert kv into tree, return HasExistedKeyError if key exists
func (t *BPTree) Insert(key uint64, value string) error {
	if t.root == nil {
		node := &Node{}
		t.root = node
		t.firstLeaf = node
		node.isLeaf = true
		node.keys = append(node.keys, key)
		node.records = append(node.records, value)
		return nil
	}
	return t.insertIntoLeaf(key, value)
}

// func (t *BPTree) Delete(key uint64) error {}
// func (t *BPTree) Find(key uint64) (string, error) {}
// func (t *BPTree) Update(key uint64, value string) error {}
