package selfbtree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	NODE_HEAD            uint16 = 0xE4E5             //节点序列化头
	KV_MAP_OFFSET_LEAF   uint32 = 2 + 1 + 2          // kv起始偏移量，叶子节点
	KV_MAP_OFFSET_BRANCH uint32 = 2 + 1 + 2 + 4096*8 // kv起始偏移量，分支节点
	CHILD_MAP_OFFSET     uint32 = 2 + 1 + 2          //子节点偏移map offset
)

type Btree interface {
	Insert(key, value Item) (bool, error)
	GET(key Item) (Item, bool, error)
	Flush() error
}

type btree struct {
	path           string       `存储路径`
	persistencer   Persistencer `持久化`
	root           *btreeNode
	capacity       uint16
	maxDegree      uint16
	reloadFromDick bool
}

type btreeNode struct {
	offset            uint64 `偏移量长度`
	len               uint32 `长度`
	holdOnMem         bool   `是否在内存中`
	bufPage           []byte
	keyOffsetMap      []uint32 `记录key的偏移量`
	kvDataOffsetStart uint32   `kv数据存储开始偏移`
	persistencer      Persistencer
	isLeaf            bool

	kv    [][]Item
	size  uint16
	child []*btreeNode
}

func (n *btreeNode) insertKv(idx int16, key Item, value Item) {
	if n.isLeaf {
		insertIdx := idx
		if idx >= 0 {
			n.kv = append(append(append([][]Item{{}}[0:0], n.kv[0:idx]...), []Item{key, value}),
				n.kv[idx:n.size]...)
		} else {
			n.kv = append(n.kv, []Item{key, value})
			insertIdx = int16(n.size)
		}
		dataLength := 4 + 4 + key.Lenth() + value.Lenth()

		kvMapBuf := n.bufPage[KV_MAP_OFFSET_LEAF : KV_MAP_OFFSET_LEAF+uint32(n.size)*4]
		copy(n.bufPage[KV_MAP_OFFSET_LEAF+uint32(insertIdx+1)*4:], kvMapBuf[insertIdx*4:])
		binary.BigEndian.PutUint32(n.bufPage[KV_MAP_OFFSET_LEAF+uint32(insertIdx)*4:], n.kvDataOffsetStart-dataLength)

		dataBuf := n.bufPage[n.kvDataOffsetStart-dataLength : n.kvDataOffsetStart]
		binary.BigEndian.PutUint32(dataBuf, dataLength)
		binary.BigEndian.PutUint32(dataBuf[4:], key.Lenth())
		copy(dataBuf[4+4:], key.Value())
		copy(dataBuf[4+4+key.Lenth():], value.Value())
		n.size++
		binary.BigEndian.PutUint16(n.bufPage[3:5], n.size)
		n.persistencer.SerializeNode(n)
		n.kvDataOffsetStart = n.kvDataOffsetStart - dataLength
	} else { //非叶子节点
		insertIdx := idx
		if idx >= 0 {
			n.kv = append(append(append([][]Item{{}}[0:0], n.kv[0:idx]...), []Item{key, nil}),
				n.kv[idx:n.size]...)
		} else {
			n.kv = append(n.kv, []Item{key, nil})
			insertIdx = int16(n.size)
		}
		dataLength := 4 + key.Lenth()

		kvMapBuf := n.bufPage[KV_MAP_OFFSET_BRANCH : KV_MAP_OFFSET_BRANCH+uint32(n.size)*4]
		copy(n.bufPage[KV_MAP_OFFSET_BRANCH+uint32(insertIdx+1)*4:], kvMapBuf[insertIdx*4:])
		binary.BigEndian.PutUint32(n.bufPage[KV_MAP_OFFSET_BRANCH+uint32(insertIdx)*4:], n.kvDataOffsetStart-dataLength)

		dataBuf := n.bufPage[n.kvDataOffsetStart-dataLength : n.kvDataOffsetStart]
		binary.BigEndian.PutUint32(dataBuf, dataLength)
		copy(dataBuf[4:], key.Value())
		n.size++
		binary.BigEndian.PutUint16(n.bufPage[3:5], n.size)
		n.persistencer.SerializeNode(n)
		n.kvDataOffsetStart = n.kvDataOffsetStart - dataLength
	}
}

func (n *btreeNode) getChild(idx int) (*btreeNode, error) {
	if idx == len(n.child) {
		fmt.Println("2")
	}
	bn := n.child[idx]
	if !bn.holdOnMem {
		cbn, err := n.persistencer.LoadNode(bn, n.child[idx].offset, 1024*64)
		if err != nil {
			return nil, err
		}
		return cbn, nil
	} else {
		return bn, nil
	}
}

//插入替换
//插入替换
func (n *btreeNode) istAndRepChilds(idx int16, left *btreeNode, right *btreeNode) {
	var child []*btreeNode
	n.child[idx] = left
	if idx == int16(n.size) {
		child = append(n.child, right)
	} else {
		child = append(n.child[0:idx+1], append([]*btreeNode{right}, n.child[idx+1:]...)...)
	}
	n.child = child
	//序列化到pageBuf
	childBuf := n.bufPage[CHILD_MAP_OFFSET : CHILD_MAP_OFFSET+4096*8]
	binary.BigEndian.PutUint64(childBuf[idx*8:], left.offset)
	if idx == int16(n.size) {
		binary.BigEndian.PutUint64(childBuf[(idx+1)*8:], right.offset)
	} else {
		copy(childBuf[(idx+2)*8:4096*8], childBuf[(idx+1)*8:4096*8])
		binary.BigEndian.PutUint64(childBuf[(idx+1)*8:], right.offset)
	}
}

func (b *btree) NewBtreeNode(orgNode *btreeNode, isLeaf bool, kv [][]Item, size uint16, childs []*btreeNode) (*btreeNode, error) {
	node := &btreeNode{
		holdOnMem:    true,
		kv:           kv,
		size:         size,
		child:        childs,
		persistencer: b.persistencer,
		isLeaf:       isLeaf,
	}

	if orgNode != nil {
		node.bufPage = orgNode.bufPage
		node.offset = orgNode.offset
	} else {
		node.bufPage = make([]byte, 1024*64)
	}

	cur := 0
	binary.BigEndian.PutUint16(node.bufPage[cur:cur+2], NODE_HEAD) //写入node head
	cur += 2
	if isLeaf {
		node.bufPage[cur] = 0x01 //最后一位为1表示是叶子节点
	} else {
		node.bufPage[cur] = 0x00
	}
	cur += 1
	if size > 4096 {
		fmt.Println("asd")
	}
	binary.BigEndian.PutUint16(node.bufPage[cur:cur+2], size) //写入元素个数
	cur += 2
	if !node.isLeaf { //只有非叶子阶段才有子树
		cnBuf := node.bufPage[cur : cur+4096*8]
		for idx, ch := range childs { //写入节点的子节点offset
			binary.BigEndian.PutUint64(cnBuf[idx*8:(idx+1)*8], ch.offset)
		}
		cur += 4096 * 8
	}
	dataOffset := uint32(cap(node.bufPage))
	if isLeaf { //叶子节点，保存键值对
		for _, itm := range kv {
			// 计算数据长度
			length := 4 + 4 + itm[0].Lenth() + itm[1].Lenth()
			// 计算数据偏移量
			dataOffset -= length
			// 保存偏移量
			binary.BigEndian.PutUint32(node.bufPage[cur:cur+4], dataOffset)
			cur += 4
			// 保存数据
			dataBuf := node.bufPage[dataOffset : dataOffset+length]
			binary.BigEndian.PutUint32(dataBuf, length)
			binary.BigEndian.PutUint32(dataBuf[4:], itm[0].Lenth())
			copy(dataBuf[8:8+itm[0].Lenth()], itm[0].Value())
			copy(dataBuf[8+itm[0].Lenth():length], itm[1].Value())
		}
	} else { //非叶子节点只保存键
		for _, itm := range kv {
			// 计算数据长度
			length := 4 + itm[0].Lenth()
			// 计算数据偏移量
			dataOffset -= length
			// 保存偏移量
			binary.BigEndian.PutUint32(node.bufPage[cur:cur+4], dataOffset)
			cur += 4
			// 保存数据
			dataBuf := node.bufPage[dataOffset : dataOffset+length]
			binary.BigEndian.PutUint32(dataBuf, length)
			copy(dataBuf[4:length], itm[0].Value())
		}
	}
	node.kvDataOffsetStart = dataOffset
	offset, err := b.persistencer.SerializeNode(node)
	if err != nil {
		return nil, err
	}
	node.offset = offset
	return node, nil
}

type Item interface {
	Less(item Item) bool
	Equal(item Item) bool
	Lenth() uint32
	Value() []byte
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

func (item *btreeItem) Lenth() uint32 {
	return uint32(len(item.key))
}

func (item *btreeItem) Value() []byte {
	return item.key
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
		binary.BigEndian.PutUint32(bs, uint32(v.(int)))
		return newItem(bs)
	}
	return nil
}

func (b *btree) Flush() error {
	return b.persistencer.Flush()
}

func (b *btree) Insert(key, value Item) (bool, error) {
	// 先看root要不要分裂
	if b.root.size == b.capacity {
		left, middle, right, err := b.split(2, b.root)
		if err != nil {
			return false, err
		}
		node, err := b.NewBtreeNode(b.root, false, [][]Item{middle}, 1, []*btreeNode{left, right})
		if err != nil {
			return false, err
		}
		b.root = node
	}
	return b.insertIntoNode(b.root, key, value), nil
}

func (b *btree) encodeNodeToBytes(node *btreeNode) {
	bs := make([]byte, 0)
	binary.BigEndian.PutUint16(bs, NODE_HEAD)
	ub := make([]byte, 8)
	binary.BigEndian.PutUint64(ub, node.offset)
	bs = append(bs, ub...)

}

func (b *btree) loadNodeFromDisk() {
}

func (b *btree) insertIntoNode(node *btreeNode, key Item, value Item) bool {
	if node.size == 0 { //空树
		node.insertKv(-1, key, value)
		return true
	} else {
		//然后插入
		var idx int16 = -1
		idx = int16(node.BinarySearchLE(key))
		if idx >= 0 && key.Equal(node.kv[idx][0]) {
			return false
		}

		if len(node.child) == 0 {
			node.insertKv(idx, key, value)
			return true
		} else {
			// 插入子节点
			childInsertIdx := idx

			if childInsertIdx < 0 {
				childInsertIdx = int16(node.size)
			}

			//校验要插入的子节点是否需要分裂
			targetNode, err := node.getChild(int(childInsertIdx))
			if err != nil {
				fmt.Println(err)
				return false
			}
			if targetNode.size == b.capacity {
				left, middle, right, err := b.split(1, targetNode)
				if err != nil {
					return false
				}
				node.istAndRepChilds(childInsertIdx, left, right)
				node.insertKv(idx, middle[0], middle[1])
				if key.Less(middle[0]) {
					return b.insertIntoNode(left, key, value)
				} else if key.Equal(middle[0]) {
					return false
				} else {
					return b.insertIntoNode(right, key, value)
				}
			}

			return b.insertIntoNode(targetNode, key, value)
		}
	}
}

func (bn *btreeNode) BinarySearchLE(key Item) int {
	start := 0
	end := int(bn.size)
	mid := int(bn.size / 2)
	kvs := bn.kv
	if bn.size == 0 {
		return -1
	}
	if key.Less(kvs[0][0]) {
		return 0
	}
	if kvs[end-1][0].Less(key) {
		return -1
	}

	for end >= start {
		if kvs[mid][0].Equal(key) {
			return mid
		} else if key.Less(kvs[mid][0]) {
			if kvs[mid-1][0].Less(key) {
				return mid
			} else {
				end = mid - 1
			}
		} else {
			start = mid + 1
		}
		mid = (start + end) / 2
	}
	return -1
}

func (b *btree) GET(key Item) (Item, bool, error) {
	value, err := b.findValueByKey(b.root, key)
	if err != nil {
		return nil, false, err
	}
	if value != nil {
		return value, true, nil
	}
	return nil, false, nil
}

func (b *btree) findValueByKey(node *btreeNode, key Item) (Item, error) {
	//for i, kv := range node.kv {
	//	if kv[0].Equal(key) {
	//		return kv[1], nil
	//	} else if key.Less(kv[0]) {
	//		if len(node.child) > 0 {
	//			cbn, err := node.getChild(i)
	//			if err != nil {
	//				return nil, err
	//			}
	//			return b.findValueByKey(cbn, key)
	//		} else {
	//			return nil, nil
	//		}
	//	}
	//}
	idx := node.BinarySearchLE(key)
	if idx >= 0 {
		tkv := node.kv[idx]
		if tkv[0].Equal(key) && node.isLeaf {
			return tkv[1], nil
		} else {
			if len(node.child) > 0 {
				cbn, err := node.getChild(idx)
				if err != nil {
					return nil, err
				}
				return b.findValueByKey(cbn, key)
			} else {
				return nil, nil
			}
		}
	}

	if len(node.child) > 0 {
		cbn, err := node.getChild(int(node.size))
		if err != nil {
			return nil, err
		}
		return b.findValueByKey(cbn, key)
	} else {
		return nil, nil
	}
}

// mode：分裂模式, 1:只分裂出一个节点，2：分裂出两个全新的节点
func (b *btree) split(mode int, parent *btreeNode) (*btreeNode, []Item, *btreeNode, error) {
	mid := parent.size - 2
	midItem := parent.kv[mid]
	isLeaf := false
	if len(parent.child) == 0 {
		isLeaf = true
	}
	var lchild, rchild []*btreeNode
	if len(parent.child) > 0 {
		if mode == 1 {
			lchild = parent.child[0 : mid+1]
		} else {
			lchild = append([]*btreeNode{}, parent.child[0:mid+1]...)
		}
		rchild = append([]*btreeNode{}, parent.child[mid+1:parent.size+1]...)
	}

	var orgNode *btreeNode = nil
	var ckv [][]Item = nil
	var leftSize uint16 = 0
	if isLeaf {
		if mode == 1 {
			orgNode = parent
			ckv = parent.kv[0 : mid+1]
		} else {
			ckv = append([][]Item{}, parent.kv[0:mid+1]...)
		}
		leftSize = mid + 1
	} else {
		if mode == 1 {
			orgNode = parent
			ckv = parent.kv[0:mid]
		} else {
			ckv = append([][]Item{}, parent.kv[0:mid]...)
		}
		leftSize = mid
	}
	left, err := b.NewBtreeNode(orgNode, isLeaf, ckv, leftSize, lchild)
	if err != nil {
		return nil, nil, nil, err
	}
	right, err := b.NewBtreeNode(nil, isLeaf, append([][]Item{}, parent.kv[mid+1:]...), 1, rchild)
	if err != nil {
		return nil, nil, nil, err
	}

	return left, midItem, right, nil
}

func NewBtree(path string, degree uint16, reloadFromDick bool) (Btree, error) {
	if degree < 10 {
		return nil, errors.New("degree is small than 2")
	}
	persistencer, err := NewSoltPersistencer(path)
	if err != nil {
		return nil, err
	}
	rootOffset, isNew := persistencer.Init()
	bt := &btree{
		path:           path,
		persistencer:   persistencer,
		maxDegree:      degree,
		capacity:       degree - 1,
		reloadFromDick: true,
	}

	if reloadFromDick && !isNew {
		root, err := persistencer.LoadNode(nil, rootOffset, 64*1024)
		if err != nil {
			return nil, err
		}
		bt.root = root
	} else {
		root, err := bt.NewBtreeNode(nil, true, make([][]Item, 0), 0, make([]*btreeNode, 0))
		if err != nil {
			return nil, err
		}
		bt.root = root
	}

	if err != nil {
		return nil, err
	}
	return bt, nil
}

func BtreeHello() {
	fmt.Println("BTREE HELLO")
}
