package selfbtree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	NODE_HEAD        uint16 = 0xE4E5         //节点序列化头
	KV_MAP_OFFSET    uint32 = 2 + 2 + 2048*8 // kv起始偏移量
	CHILD_MAP_OFFSET uint32 = 2 + 2          //子节点偏移map offset
)

type Btree interface {
	Insert(key, value Item) bool
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

	kv    [][]Item
	size  uint16
	child []*btreeNode
}

func (n *btreeNode) insertKv(idx int16, key Item, value Item) {
	insertIdx := idx
	if idx >= 0 {
		n.kv = append(append(append([][]Item{{}}[0:0], n.kv[0:idx]...), []Item{key, value}),
			n.kv[idx:n.size]...)
	} else {
		n.kv = append(n.kv, []Item{key, value})
		insertIdx = int16(n.size)
	}
	dataLength := 4 + 4 + key.Lenth() + value.Lenth()

	kvMapBuf := n.bufPage[KV_MAP_OFFSET : KV_MAP_OFFSET+uint32(n.size)*4]
	copy(n.bufPage[KV_MAP_OFFSET+uint32(insertIdx+1)*4:], kvMapBuf[insertIdx*4:])
	binary.BigEndian.PutUint32(n.bufPage[KV_MAP_OFFSET+uint32(insertIdx)*4:], n.kvDataOffsetStart-dataLength)

	if n.kvDataOffsetStart < KV_MAP_OFFSET {
		fmt.Println(n.kvDataOffsetStart)
	}
	dataBuf := n.bufPage[n.kvDataOffsetStart-dataLength : n.kvDataOffsetStart]
	binary.BigEndian.PutUint32(dataBuf, dataLength)
	binary.BigEndian.PutUint32(dataBuf[4:], key.Lenth())
	copy(dataBuf[4+4:], key.Value())
	copy(dataBuf[4+4+key.Lenth():], value.Value())
	n.size++
	binary.BigEndian.PutUint16(n.bufPage[2:4], n.size)
	n.persistencer.SerializeNode(n)
	n.kvDataOffsetStart = n.kvDataOffsetStart - dataLength
}

func (n *btreeNode) getChild(idx int) (*btreeNode, error) {
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
func (n *btreeNode) istAndRepChilds(idx int16, left *btreeNode, right *btreeNode) {
	if len(n.child) == 192 {
		fmt.Printf("idx:=%d,len(n.child)=%d", idx, len(n.child))
	}
	var child []*btreeNode
	n.child[idx] = left
	if idx == int16(n.size) {
		child = append(n.child, right)
	} else {
		child = append(n.child[0:idx+1], append([]*btreeNode{right}, n.child[idx+1:]...)...)
	}
	n.child = child
	//序列化到pageBuf
	childBuf := n.bufPage[CHILD_MAP_OFFSET : CHILD_MAP_OFFSET+2048*8]
	binary.BigEndian.PutUint64(childBuf[idx*8:], left.offset)
	if idx == int16(n.size) {
		binary.BigEndian.PutUint64(childBuf[(idx+1)*8:], right.offset)
	} else {
		copy(childBuf[(idx+2)*8:2048*8], childBuf[(idx+1)*8:2048*8])
		binary.BigEndian.PutUint64(childBuf[(idx+1)*8:], right.offset)
	}
}

func (b *btree) NewBtreeNode(orgNode *btreeNode, kv [][]Item, size uint16, childs []*btreeNode) (*btreeNode, error) {
	node := &btreeNode{
		holdOnMem:    true,
		kv:           kv,
		size:         size,
		child:        childs,
		persistencer: b.persistencer,
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
	binary.BigEndian.PutUint16(node.bufPage[cur:cur+2], size) //写入元素个数
	cur += 2
	cnBuf := node.bufPage[cur : cur+2048*8]
	cur += 2048 * 8
	for idx, ch := range childs { //写入节点的子节点offset
		binary.BigEndian.PutUint64(cnBuf[idx*8:(idx+1)*8], ch.offset)
	}
	dataOffset := uint32(cap(node.bufPage))
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
		copy(dataBuf[8:length], append(append([]byte{}, itm[0].Value()...), itm[1].Value()...))
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

func (b *btree) Insert(key, value Item) bool {
	state, err := b.insertIntoNode(b.root, key, value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if state == 1 {
		return true
	} else if state == 3 { //空间不足
		right, err := b.NewBtreeNode(nil, make([][]Item, 0), 0, make([]*btreeNode, 0))
		if err != nil {
			return false
		}
		root, err := b.NewBtreeNode(nil, append([][]Item{}, []Item{key, value}),
			1, []*btreeNode{b.root, right})
		if err != nil {
			return false
		}
		b.root = root
		b.persistencer.UpdateRoot(b.root)
		return true
	} else {
		return false
	}
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

// 1:插入成功，2:键重复，3:空间不足，4：其他
func (b *btree) insertIntoNode(node *btreeNode, key Item, value Item) (uint8, error) {
	if node.size == 0 { //空树
		node.insertKv(-1, key, value)
		return 1, nil
	} else { //没有子树时
		//然后插入
		var idx int16 = -1
		for i, kv := range node.kv {
			if key.Less(kv[0]) {
				idx = int16(i)
				break
			} else if key.Equal(kv[0]) {
				return 2, nil
			}
		}
		if len(node.child) == 0 {
			if node.size == b.capacity {
				return 3, nil
			}
			node.insertKv(idx, key, value)
			return 1, nil
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
				return 4, err
			}

			state, err := b.insertIntoNode(targetNode, key, value)
			if state == 3 {
				if node.size == b.capacity {
					return 3, nil
				}
				fmt.Println("split...")
				right, err := b.NewBtreeNode(nil, make([][]Item, 0), 0, make([]*btreeNode, 0))
				if err != nil {
					return 4, err
				}
				node.istAndRepChilds(childInsertIdx, targetNode, right)
				node.insertKv(idx, key, value)
				return 1, nil
			} else {
				return state, err
			}
		}
	}
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
	for i, kv := range node.kv {
		if kv[0].Equal(key) {
			return kv[1], nil
		} else if key.Less(kv[0]) {
			if len(node.child) > 0 {
				cbn, err := node.getChild(i)
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
	mid := parent.size / 2
	midItem := parent.kv[mid]

	var lchild, rchild []*btreeNode
	if len(parent.child) > 0 {
		lchild = parent.child[0 : mid+1]
		rchild = parent.child[0 : mid+1]
	}

	var orgNode *btreeNode = nil
	if mode == 1 {
		orgNode = parent
	}
	left, err := b.NewBtreeNode(orgNode, parent.kv[0:mid], mid, lchild)
	if err != nil {
		return nil, nil, nil, err
	}
	right, err := b.NewBtreeNode(nil, parent.kv[mid+1:], mid-1, rchild)
	if err != nil {
		return nil, nil, nil, err
	}

	return left, midItem, right, nil
}

func NewBtree(path string, degree uint16, reloadFromDick bool) (Btree, error) {
	if degree < 2 {
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
		root, err := bt.NewBtreeNode(nil, make([][]Item, 0), 0, make([]*btreeNode, 0))
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
