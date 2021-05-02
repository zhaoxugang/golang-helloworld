package selfbtree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const NODE_HEAD = 0xE4E5 //节点序列化头

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
	return b.insertIntoNode(b.root, key, value)
}

func (b *btree) encodeNodeToBytes(node *btreeNode) {
	bs := make([]byte, 0)
	bs = append(bs, NODE_HEAD)
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
		if len(node.child) <= 0 {
			if idx >= 0 {
				node.kv = append(append(append([][]Item{{}}[0:0], node.kv[0:idx]...), []Item{key, value}),
					node.kv[idx:node.size]...)
				node.size++
			} else {
				node.kv = append(node.kv, []Item{key, value})
				node.size++
			}
			if node.size > b.capacity {
				//分裂
				middle := node.size / 2
				newNode := &btreeNode{
					kv:     append([][]Item{{}}[0:0], node.kv[middle+1:node.size]...),
					size:   node.size - middle - 1,
					parent: node.parent,
				}
				b.insertNode(newNode.parent, node, newNode, node.kv[middle])
				node.kv = node.kv[0:middle]
				node.size = middle
			}
			return true
		} else {
			// 插入子节点
			if idx >= 0 {
				return b.insertIntoNode(node.child[idx], key, value)
			} else {
				return b.insertIntoNode(node.child[node.size], key, value)
			}
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

func (b *btree) insertNode(parent *btreeNode, left, right *btreeNode, items []Item) {
	key := items[0]
	value := items[1]
	if parent == nil {
		parent = &btreeNode{
			kv:    make([][]Item, 0),
			size:  0,
			child: make([]*btreeNode, 0),
		}
		parent.child = append(parent.child, left, right)
		parent.kv = append(parent.kv, items)
		b.root = parent
		left.parent = parent
		right.parent = parent
		parent.size++
	} else {
		// 先插入
		idx := -1
		for i, kv := range parent.kv {
			if key.Less(kv[0]) {
				idx = i
				break
			}
		}
		if idx >= 0 {
			parent.kv = append(append(append([][]Item{{}}[0:0], parent.kv[0:idx]...), []Item{key, value}),
				parent.kv[idx:parent.size]...)
			parent.child = append(append(append([]*btreeNode{}[0:0], parent.child[0:idx+1]...), right),
				parent.child[idx+1:parent.size+1]...)
			right.parent = parent
			parent.size++
		} else {
			parent.kv = append(parent.kv, []Item{key, value})
			parent.child = append(parent.child, right)
			right.parent = parent
			parent.size++
		}
		//分裂
		if parent.size > b.capacity {
			middle := parent.size / 2
			newNode := &btreeNode{
				kv:     append([][]Item{{}}[0:0], parent.kv[middle+1:parent.size]...),
				size:   parent.size - middle - 1,
				parent: parent.parent,
				child:  append([]*btreeNode{}[0:0], parent.child[middle+1:parent.size+1]...),
			}
			for _, bnode := range newNode.child {
				bnode.parent = newNode
			}
			b.insertNode(parent.parent, parent, newNode, parent.kv[middle])
			parent.kv = parent.kv[0:middle]
			parent.child = parent.child[0 : middle+1]
			parent.size = middle
		}
	}
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
