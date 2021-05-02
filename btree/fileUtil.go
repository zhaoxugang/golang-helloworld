package selfbtree

import (
	"errors"
	"os"
)

func readFile(f *os.File, offset uint64, buf []byte) error {
	_, err := f.ReadAt(buf, int64(offset))
	if err != nil {
		return errors.New("读取文件失败")
	}
	return nil
}

func writeFile(f *os.File, offset uint64, buf []byte) error {
	_, err := f.WriteAt(buf, int64(offset))
	if err != nil {
		return errors.New("写入文件失败")
	}
	return nil
}
