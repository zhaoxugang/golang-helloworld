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
	LoadNode(offset uint64, size uint32) (*btreeNode, error)
	SerializeNode(node *btreeNode) error
}

//分槽页
type SoltPersistencer struct {
	magic   []byte
	idxFile *os.File
}

func NewSoltPersistencer(path string) (*SoltPersistencer, error) {
	idxFile, err := os.Open(path)
	if err != nil {
		return nil, errors.New("打开文件失败")
	}
	defer idxFile.Close()
	return &SoltPersistencer{
		magic:   []byte{0xEF, 0xEF},
		idxFile: idxFile,
	}, nil
}

func (p *SoltPersistencer) LoadNode(offset uint64, size uint32) (*btreeNode, error) {
	buf := make([]byte, size)
	err := readFile(p.idxFile, offset, buf)
	if err != nil {
		return nil, err
	}
	bn, err := p.decoder(buf)
	if err != nil {
		return nil, err
	}
	return bn, nil
}

func (p *SoltPersistencer) SerializeNode(node *btreeNode) error {
	err := writeFile(p.idxFile, node.offset, node.bufPage)
	return err
}

func (p *SoltPersistencer) decoder(buf []byte) (*btreeNode, error) {
	bn := &btreeNode{
		holdOnMem: true,
		child:     make([]*btreeNode, 128),
		bufPage:   buf,
	}
	cur := 0
	cp := bytes.Compare(p.magic[cur:2], buf[cur:2])
	if cp != 0 {
		return nil, errors.New("magic异常")
	}
	cur += 2

	//节点元素数
	size := int16(binary.LittleEndian.Uint16(buf[cur:2])) //unint32类型
	bn.size = size
	cur += 2

	// 节点子节点引用
	childOffets := buf[cur : cur+128*8] // 128个子节点
	for i := 0; i < len(childOffets); i += 8 {
		cbn := &btreeNode{
			holdOnMem: false,
			offset:    binary.LittleEndian.Uint64(childOffets[i*8 : (i+1)*8]),
		}
		bn.child[i] = cbn
		cbn.parent = bn
	}
	cur += 128 * 8

	// kv存储直接放在buf中
	for i := 0; i < int(size); i++ {
		offset := binary.BigEndian.Uint32(buf[cur+i*4 : cur+(i+1)*4])
		len := binary.BigEndian.Uint32(buf[offset : offset+4])
		kvbuf := buf[offset+4 : offset+4+len]
		dataOffset := binary.BigEndian.Uint32(kvbuf[0:4])
		key := &btreeItem{key: kvbuf[0:dataOffset]}
		value := &btreeItem{key: kvbuf[dataOffset:]}
		bn.keyOffsetMap[i] = offset
		bn.kv = append(bn.kv, []Item{key, value})
	}

	return bn, nil
}
