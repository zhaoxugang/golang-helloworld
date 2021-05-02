package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

type Employee struct {
	ID        int
	Name      string
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	ManagerId int
}

var dilbert Employee

func main32() {
	//dilbert.Salary -= 5000
	//position := &dilbert.Position
	//*position = "Seniot " + *position
	//fmt.Printf("%s", *position)
	//
	//fmt.Print(dilbert.Address == "")
	//var employeeOfTheMonth *Employee = &dilbert
	//employeeOfTheMonth.Position += "(proactive team player)"
	//(*employeeOfTheMonth).Position = "(proactive team player" //与上面的用法等价
	//
	//a := getEmployed().Salary
	//a = 0
	//fmt.Printf("%d", a)

	var testStruct TestStruct = TestStruct{hello: hello}
	fmt.Printf(testStruct.hello("zhaoxugang"))
	var a uint64
	fmt.Printf("\n%d\n", a)
	var s []string
	fmt.Println(s, len(s), cap(s))
	fmt.Println(s == nil)

	s = append(s, "Hello")
	s = append(s, "world")
	fmt.Println(strings.Join(s, ", "))
	var b *int
	b = nil
	var c interface{}
	var d interface{}
	c = b
	d = c
	fmt.Println("------1-------")
	fmt.Println(c == nil)
	fmt.Println(b == nil)
	fmt.Println(d == nil)
	fmt.Println(c == b)
	fmt.Println("------2-------")

	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)

	var p *int
	var i interface{}
	fmt.Println(p == nil)
	fmt.Println(i == nil)
	i = p
	fmt.Println(i == nil)

	var q Column
	fmt.Println(unsafe.Sizeof(q))
	fmt.Println(unsafe.Alignof(q.b))
	fmt.Println(unsafe.Alignof(q.length))
	fmt.Println("=============")
	var x [10000000]struct{}
	x[0] = struct{}{}
	fmt.Println(len(x))
	fmt.Println(unsafe.Sizeof(x))
	var y = make([]struct{}, 1000000)
	fmt.Println(unsafe.Sizeof(y))
	fmt.Println(reflect.TypeOf(nil))
	fmt.Print("======\n")
	arr1 := make([]int, 4, 32)
	fmt.Printf("=====%d+++", arr1[2])
}

func Foo() error {
	var err *os.PathError
	return err
}

type Column struct {
	b          byte
	length     int
	nullBitmap []byte // bit 0 is null, 1 is not null
	offsets    []int64
	data       []byte
	elemBuf    []byte
}

func HashGroupKey(buf [][]byte) {
	fmt.Printf("%d", len(buf))
}

type TestStruct struct {
	hello func(str string) string
}

func hello(str string) string {
	s := fmt.Sprintf("%s %s", str, "hello")
	return s
}

func getEmployed() Employee {
	return dilbert
}
