package main

import (
	"bytes"
	"fmt"
	"unsafe"
)

func main20() {
	var buffer bytes.Buffer
	fmt.Println(unsafe.Sizeof(buffer))
	buffer.WriteString("Hello")
	fmt.Println(buffer)

	var a []int
	a = nil
	fmt.Println(len(a))
}
