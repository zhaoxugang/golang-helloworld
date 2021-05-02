package selfbtree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const NODE_HEAD uint16 = 0xE4E5 //节点序列化头

type Btree interface {
	Insert(key, value Item) bool
	GET(key Item) (Item, bool)
}

type btree struct {
	path         string       `存储路径`
	persistencer Persistencer `持久化`
	root         *btreeNode
	capacity     int16
	maxDegree    int16
}

type btreeNode struct {
	offset       uint64 `偏移量长度`
	len          uint32 `长度`
	holdOnMem    bool   `是否在内存中`
	bufPage      []byte
	keyOffsetMap []uint32 `记录key的偏移量`

	kv     [][]Item
	size   int16
	parent *btreeNode
	child  []*btreeNode
}

type Item interface {
	Less(item Item) bool
	Equal(item Item) bool
}

type btreeItem struct {
	key []byte
}

func (item *btreeItem) Less(other Item) bool {
	return bytes.Compare(item.key, other.(*btreeItem).key) < 0
}

func (item *btreeItem) Equal(other Item) bool {
	return bytes.Compare(item.key, other.(*btreeItem).key) == 0
}

func newItem(bs []byte) Item {
	return &btreeItem{
		key: bs,
	}
}

func Encode(v interface{}) Item {
	switch v.(type) {
	case int:
		bs := make([]byte, 4)
		for i := 3; i >= 0; i-- {
			iv := v.(int)
			bs[i] = byte(iv & 0xff)
			v = iv >> 8
		}
		return newItem(bs)
	}
	return nil
}

func (b *btree) Insert(key, value Item) bool {
	// 先看root要不要分裂
	if b.root.size == b.maxDegree {
		left, middle, right := b.split(b.root)
		b.root = &btreeNode{}
		b.root.kv = append(b.root.kv, middle)
		b.root.child = append(b.root.child, left, right)
		b.root.size = 1
		left.parent = b.root
		right.parent = b.root
	}
	return b.insertIntoNode(b.root, key, value)
}

func (b *btree) encodeNodeToBytes(node *btreeNode) {
	bs := make([]byte, 0)
	binary.BigEndian.PutUint16(bs, NODE_HEAD)
	ub := make([]byte, 8)
	binary.LittleEndian.PutUint64(ub, node.offset)
	bs = append(bs, ub...)

}

func (b *btree) loadNodeFromDisk() {
}

func (b *btree) insertIntoNode(node *btreeNode, key Item, value Item) bool {
	if node.size == 0 { //空树
		node.kv = append(node.kv, []Item{key, value})
		node.size++
		return true
	} else { //没有子树时
		//然后插入
		idx := -1
		for i, kv := range node.kv {
			if key.Less(kv[0]) {
				idx = i
				break
			} else if key.Equal(kv[0]) {
				return false
			}
		}
		if len(node.child) == 0 {
			if idx >= 0 {
				node.kv = append(append(append([][]Item{{}}[0:0], node.kv[0:idx]...), []Item{key, value}),
					node.kv[idx:node.size]...)
				node.size++
			} else {
				node.kv = append(node.kv, []Item{key, value})
				node.size++
			}
			return true
		} else {
			// 插入子节点
			if idx < 0 {
				idx = int(node.size)
			}

			//校验要插入的子节点是否需要分裂
			targetNode := node.child[idx]
			if targetNode.size == b.maxDegree {
				left, middle, right := b.split(targetNode)
				targetNode = &btreeNode{
					size:   1,
					parent: node,
				}
				targetNode.kv = append(targetNode.kv, middle)
				targetNode.child = append(targetNode.child, left, right)

				node.child[idx] = targetNode
				left.parent = targetNode
				right.parent = targetNode
			}

			return b.insertIntoNode(targetNode, key, value)
		}
	}
}

func (b *btree) GET(key Item) (Item, bool) {
	value := b.findValueByKey(b.root, key)
	if value != nil {
		return value, true
	}
	return nil, false
}

func (b *btree) findValueByKey(node *btreeNode, key Item) Item {
	for i, kv := range node.kv {
		if kv[0].Equal(key) {
			return kv[1]
		} else if key.Less(kv[0]) {
			if len(node.child) > 0 {
				return b.findValueByKey(node.child[i], key)
			} else {
				return nil
			}
		}
	}
	if len(node.child) > 0 {
		return b.findValueByKey(node.child[node.size], key)
	} else {
		return nil
	}
}

func (b *btree) split(parent *btreeNode) (*btreeNode, []Item, *btreeNode) {
	mid := parent.size / 2
	midItem := parent.kv[mid]
	left := &btreeNode{
		kv:   parent.kv[0:mid],
		size: mid,
	}
	right := &btreeNode{
		kv:   parent.kv[mid+1:],
		size: mid - 1,
	}

	if len(parent.child) > 0 {
		left.child = parent.child[0 : mid+1]
		right.child = parent.child[0 : mid+1]
	}
	return left, midItem, right
}

func NewBtree(path string, degree int16) (Btree, error) {
	if degree < 2 {
		return nil, errors.New("degree is small than 2")
	}
	root := &btreeNode{
		kv:    make([][]Item, 0),
		size:  0,
		child: make([]*btreeNode, 0),
	}
	persistencer, err := NewSoltPersistencer(path)
	if err != nil {
		return nil, err
	}
	return &btree{
		path:         path,
		persistencer: persistencer,
		maxDegree:    degree,
		capacity:     degree - 1,
		root:         root,
	}, nil
}

func BtreeHello() {
	fmt.Println("BTREE HELLO")
}
