package selfbtree

import (
	"errors"
	"os"
)

func readFile(f *os.File, offset uint64, buf []byte) error {
	_, err := f.ReadAt(buf, int64(offset))
	if err != nil {
		return err
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

func appendFile(f *os.File, buf []byte) (int64, error) {
	offset, err := f.Seek(0, 2)
	if err != nil {
		return -1, err
	}
	_, err = f.WriteAt(buf, int64(offset))
	if err != nil {
		return -1, errors.New("写入文件失败")
	}
	return offset, nil
}
