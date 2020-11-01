package memory

import (
	"errors"
	"sync"
)

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

func (t *BPTree) Delete(key uint64) error {
	//删除一个节点，先从树根往下查找这个节点
	t.root.remove(key, t)
	return errors.New("success")
}
func (n *Node) remove(key uint64, tree *BPTree) error {
	//如果是叶子节点
	if n.isLeaf {
		//叶节点中不存在关键字，直接返回未找到信息
		index := n.containsKey(key)
		if index == -1 {
			return errors.New("No key founded")
		}
		//如果既是叶节点又是根节点，
		if n.parent == nil {
			n.removeLeafKey(index)
		} else {
			if n.canRemoveDirectly() { //能在叶子节点中直接删除
				n.removeLeafKey(index)
			} else {
				if n.canborrow(n.prev) { //左兄弟能够借出
					n.removeLeafKey(index)
					n.borrowFromLeftLeaf()
				} else if n.canborrow(n.next) { //右兄弟能够借出
					n.removeLeafKey(index)
					n.borrowFromRightLeaf()
				} else { //需要将两个叶子节点合并
					if n.canMerge(n.prev) { //和左兄弟合并
						n.removeLeafKey(index)
						n.addPreNode(n.prev)
						n.prev.parent = nil
						n.prev.keys = nil
						n.prev.records = nil
						curIndex := n.getParentIndex()
						//移除父节点中的key索引的值
						n.parent.removeInsideKey(uint64(curIndex))
						//删除父节点中指向左兄弟的指针
						n.parent.children = append(n.parent.children[0:curIndex], n.parent.children[curIndex+1:]...)
						//更新链表
						if n.prev.prev != nil {
							temp := n.prev
							temp.prev.next = n
							n.prev = temp.prev
							temp.prev = nil
							temp.next = nil
						} else {
							//设置头节点
						}

					} else if n.canMerge(n.next) { //和右兄弟合并
						n.removeLeafKey(index)
						n.addNextNode(n.next)
						n.next.parent = nil
						n.next.keys = nil
						n.records = nil
						curIndex := n.next.getParentIndex()
						//移除父节点中的key索引的值
						n.parent.removeInsideKey(uint64(curIndex))
						//删除父节点指向右兄弟的指针
						n.parent.children = append(n.parent.children[0:curIndex+1], n.parent.children[curIndex+2:]...)
						//更新链表
						if n.next.next != nil {
							temp := n.next
							temp.next.prev = n
							n.next = temp.next
							temp.prev = nil
							temp.next = nil
						}
					} else {

					}
				}

			}
			n.parent.updateRemove(tree)
		}
	} else {
		//递归查询
		if key < n.keys[0] {
			//沿着第一个子节点进行删除
			return n.children[0].remove(key, tree)
		} else if key >= n.keys[len(n.keys)-1] {
			//沿着最后一个子节点进行删除
			return n.children[len(n.children)-1].remove(key, tree)
		} else {
			for i := 1; i < len(n.keys); i++ {
				if key < n.keys[i] {
					return n.children[i].remove(key, tree)
				}
			}
		}
	}
	return errors.New("success")
}

//删除节点后的内部节点进行跟新
func (n *Node) updateRemove(tree *BPTree) {
	//判断该节点是否满足B+树性质，即节点数>order/2
	if len(n.children) <= order/2 {
		if n.parent == nil {
			//当前节点为根
			if len(n.children) >= 2 {
				//根节点至少两个孩子节点，符合B+树性质，停止修复，返回
				return
			} else {
				//直接将子节点作为根节点
				node := n.children[0]
				tree.root = node
				tree.root.parent = nil
				n.children = nil
				n.keys = nil
			}
		} else { //中间节点修复
			curIndex := n.parent.getParentIndex() + 1
			preIndex := curIndex - 1
			nextIndex := curIndex + 1
			var preNode, nextNode *Node = nil, nil
			if preIndex >= 0 {
				preNode = n.parent.children[preIndex]
			}
			if nextIndex < len(n.parent.children) {
				nextNode = n.parent.children[nextIndex]
			}
			if n.canborrow(preNode) {
				n.borrowNodePrevious(preNode)
			} else if n.canborrow(nextNode) {
				n.borrowNodeNext(nextNode)
			} else {
				//将两个节点合并
				if n.canMerge(preNode) {
					n.addPreNode(preNode)
					preNode.parent = nil
					preNode.keys = nil
					curIndexkey := n.getParentIndex()
					//删除父节点中对当前节点的key索引
					n.parent.keys = append(n.parent.keys[0:curIndexkey], n.parent.keys[curIndexkey+1:]...)
					//删除父节点中对前驱节点的索引
					n.parent.children = append(n.parent.children[0:curIndexkey], n.parent.children[curIndexkey+1:]...)
				} else if n.canMerge(nextNode) {
					n.addNextNode(nextNode)
					nextNode.parent = nil
					nextNode.keys = nil
					curIndexkey := nextNode.getParentIndex()
					n.parent.keys = append(n.parent.keys[0:curIndexkey], n.parent.keys[curIndexkey+1:]...)
					n.parent.children = append(n.parent.children[0:curIndexkey+1], n.parent.children[curIndexkey+2:]...)
				}
			}
			n.parent.updateRemove(tree)
		}
	}
}

//向右兄弟借值
func (n *Node) borrowNodeNext(nextNode *Node) {
	parIndex := n.getParentIndex() + 1
	downkey := n.parent.keys[parIndex]
	n.keys = append(n.keys, downkey) //将key下移
	//将nextNode的key上移
	n.parent.keys[parIndex] = nextNode.keys[0]
	//移除nextNode的第一个key
	nextNode.keys = nextNode.keys[1:]
	//将nextNode的child节点移到到当前节点
	n.children = append(n.children, nextNode.children[0])
	n.children[len(n.children)-1].parent = n
	nextNode.children = nextNode.children[1:]
}

//向左兄弟借值
func (n *Node) borrowNodePrevious(preNode *Node) {
	parIndex := n.getParentIndex()
	downkey := n.parent.keys[parIndex]
	n.keys = append([]uint64{downkey}, n.keys...) //先将key下移
	//将preNode的key上移动
	n.parent.keys[parIndex] = preNode.keys[len(preNode.keys)-1]
	preNode.keys = preNode.keys[:len(preNode.keys)-1]
	//将preNode的child移到当前节点
	n.children = append([]*Node{preNode.children[len(preNode.children)-1]}, n.children...)
	//修改转移过来的child的父节点
	n.children[0].parent = n
	preNode.children = preNode.children[:len(preNode.children)-1]
}

//当前节点与右兄弟合并
func (n *Node) addNextNode(nextnode *Node) {
	if !nextnode.isLeaf {
		index := nextnode.getParentIndex()
		n.keys = append(n.keys, nextnode.parent.keys[index])
	}
	n.keys = append(n.keys, nextnode.keys...)
	if !nextnode.isLeaf {
		//设置父节点
		for i := 0; i < len(nextnode.children); i++ {
			nextnode.children[i].parent = n
		}
		n.children = append(n.children, nextnode.children...)
	}
}

//当前节点与前一个节点合并
func (n *Node) addPreNode(prenode *Node) {
	if !prenode.isLeaf {
		index := n.getParentIndex()
		tempkey := []uint64{n.parent.keys[index]}
		n.keys = append(tempkey, n.keys...)
	}
	tempkey := prenode.keys
	n.keys = append(tempkey, n.keys...)
	if !prenode.isLeaf {
		//设置父节点
		for i := 0; i < len(prenode.children); i++ {
			prenode.children[i].parent = n
		}
		n.children = append(prenode.children, n.children...)
	}

}
func (n *Node) canMerge(merge *Node) bool {
	if merge != nil && merge.parent == n.parent && len(merge.keys) <= order/2 {
		return true
	}
	return false
}

//从右兄弟借值
func (n *Node) borrowFromRightLeaf() {
	// 从右借第一个过来,加到当前节点的最后面
	keyborrow := n.next.keys[0]
	valborrow := n.next.records[0]
	n.next.removeLeafKey(0)
	n.keys = append(n.keys, keyborrow)
	n.records = append(n.records, valborrow)
	// 找到当前节点的后继节点在父节点中的索引
	index := n.next.getParentIndex()
	//修改父节点key中的索引值
	n.parent.keys[index] = n.next.keys[0]
}

//从左兄弟借值
func (n *Node) borrowFromLeftLeaf() {
	size := len(n.prev.keys)
	keyborrow := n.prev.keys[size-1]
	valborrow := n.prev.records[size-1]
	n.prev.removeLeafKey(uint64(size - 1))

	tempkey := []uint64{keyborrow}
	tempkey = append(tempkey, n.keys...)
	n.keys = tempkey
	tempval := []string{valborrow}
	tempval = append(tempval, n.records...)
	n.records = tempval
	// 找到当前节点在父节点中的索引位置
	index := n.getParentIndex()
	//修改父节点key中的索引值
	n.parent.keys[index] = keyborrow
}

//找到当前节点在父节点中的entries
func (n *Node) getParentIndex() int {
	for i, child := range n.parent.children {
		if child == n {
			return i - 1
		}
	}
	return -1
}

func (n *Node) removeLeafKey(index uint64) {
	n.keys = append(n.keys[:index], n.keys[index+1:]...)
	n.records = append(n.records[:index], n.records[index+1:]...)
}

func (n *Node) removeInsideKey(index uint64) {
	n.keys = append(n.keys[:index], n.keys[index+1:]...)
}

//判断某个节点（同属于一个父节点）是否有多余的值可以借出
func (n *Node) canborrow(borrow *Node) bool {
	if borrow != nil && len(borrow.keys) > order/2 && borrow.parent == n.parent {
		return true
	}
	return false
}

func (n *Node) canRemoveDirectly() bool {
	if len(n.keys) > order/2 {
		return true
	}
	return false
}

//在node中的关键字切片中进行二分法查找，返回元素下标，若没找着返回-1
func (n *Node) containsKey(key uint64) uint64 {
	var low, high uint64 = 0, uint64(len(n.keys) - 1)
	var mid uint64
	for low <= high {
		mid = uint64(low+high) / 2
		if key == n.keys[mid] {
			return mid
		}
		if key < n.keys[mid] {
			high = mid - 1
		}
		if key > n.keys[mid] {
			low = mid + 1
		}
	}
	return -1
}
