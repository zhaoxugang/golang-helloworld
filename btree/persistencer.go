package selfbtree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

/**
持久化
*/
type Persistencer interface {
	LoadNode(bn *btreeNode, offset uint64, size uint32) (*btreeNode, error)
	SerializeNode(node *btreeNode) (uint64, error)
	Flush() error
	UpdateRoot(bn *btreeNode) bool
}

//分槽页
type SoltPersistencer struct {
	magic   []byte
	idxFile *os.File
}

func NewSoltPersistencer(path string) (*SoltPersistencer, error) {
	idxFile, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		return nil, errors.New("打开文件失败")
	}
	defer func() {
		if err := recover(); err != nil {
			idxFile.Close()
		}
	}()
	return &SoltPersistencer{
		magic:   []byte{0xEF, 0xEF},
		idxFile: idxFile,
	}, nil
}

func (p *SoltPersistencer) Flush() error {
	return p.idxFile.Sync()
}

func (p *SoltPersistencer) LoadNode(bn *btreeNode, offset uint64, size uint32) (*btreeNode, error) {
	buf := make([]byte, size)
	err := readFile(p.idxFile, offset, buf)
	if err != nil {
		return nil, err
	}
	if bn == nil {
		bn = &btreeNode{
			offset:       offset,
			holdOnMem:    true,
			child:        make([]*btreeNode, 0),
			bufPage:      buf,
			keyOffsetMap: make([]uint32, 3072),
			persistencer: p,
		}
	} else {
		bn.holdOnMem = true
		bn.child = make([]*btreeNode, 0)
		bn.bufPage = buf
		bn.keyOffsetMap = make([]uint32, 3072)
		bn.persistencer = p
	}
	err = p.decoder(bn, buf)
	if err != nil {
		return nil, err
	}
	return bn, nil
}

func (p *SoltPersistencer) SerializeNode(node *btreeNode) (uint64, error) {
	if node.offset > 0 {
		err := writeFile(p.idxFile, node.offset, node.bufPage)
		return node.offset, err
	} else {
		offset, err := appendFile(p.idxFile, node.bufPage)
		if err != nil {
			return 0, err
		}
		return uint64(offset), err
	}
}

func (p *SoltPersistencer) decoder(bn *btreeNode, buf []byte) error {
	cur := 0
	nodeMagic := binary.BigEndian.Uint16(buf[cur:])
	if NODE_HEAD != nodeMagic {
		return errors.New("magic异常")
	}
	cur += 2
	nodeMata := buf[cur]
	if nodeMata&0x01 > 0 {
		bn.isLeaf = true
	}
	cur += 1
	//节点元素数
	size := binary.BigEndian.Uint16(buf[cur:]) //unint16类型
	bn.size = size
	cur += 2
	if !bn.isLeaf {
		// 节点子节点引用
		childOffets := buf[cur : cur+3072*8] // 3072个子节点
		for i := 0; i <= int(size); i += 1 {
			coff := binary.BigEndian.Uint64(childOffets[i*8 : (i+1)*8])
			cbn := &btreeNode{
				holdOnMem:    false,
				offset:       coff,
				persistencer: p,
			}
			bn.child = append(bn.child, cbn)
		}
		cur += 3072 * 8
	}

	kvDataOffsetStart := uint32(len(buf))
	// kv存储直接放在buf中
	if bn.isLeaf {
		for i := 0; i < int(size); i++ {
			offset := binary.BigEndian.Uint32(buf[cur+i*4 : cur+(i+1)*4])
			len := binary.BigEndian.Uint32(buf[offset : offset+4])
			kvbuf := buf[offset+4 : offset+len]
			keyLen := binary.BigEndian.Uint32(kvbuf[0:4])
			key := &btreeItem{key: kvbuf[4 : 4+keyLen]}
			value := &btreeItem{key: kvbuf[4+keyLen:]}
			bn.keyOffsetMap[i] = offset
			if offset < kvDataOffsetStart {
				kvDataOffsetStart = offset
			}
			bn.kv = append(bn.kv, []Item{key, value})
		}
	} else {
		for i := 0; i < int(size); i++ {
			offset := binary.BigEndian.Uint32(buf[cur+i*4 : cur+(i+1)*4])
			len := binary.BigEndian.Uint32(buf[offset : offset+4])
			kvbuf := buf[offset+4 : offset+len]
			key := &btreeItem{key: kvbuf}
			bn.keyOffsetMap[i] = offset
			if offset < kvDataOffsetStart {
				kvDataOffsetStart = offset
			}
			bn.kv = append(bn.kv, []Item{key, nil})
		}
	}
	bn.kvDataOffsetStart = kvDataOffsetStart

	return nil
}

func (p *SoltPersistencer) UpdateRoot(bn *btreeNode) bool {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, bn.offset)
	err := writeFile(p.idxFile, 2, buf)
	if err != nil {
		return false
	}
	return true
}

func (p *SoltPersistencer) Init() (uint64, bool) {
	buf := make([]byte, 8)
	readFile(p.idxFile, 0, buf[0:2])

	if bytes.Compare(buf[0:2], p.magic) != 0 {
		writeFile(p.idxFile, 0, p.magic) //写入magic
		binary.BigEndian.PutUint64(buf, 10)
		writeFile(p.idxFile, 2, buf)
		return 10, true
	}
	readFile(p.idxFile, 2, buf)
	return binary.BigEndian.Uint64(buf), false
}
