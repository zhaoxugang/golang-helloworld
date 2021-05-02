package main

import (
	"fmt"
)

func main333() {
	var a [3]int
	fmt.Println(a[0])
	fmt.Println(a[len(a)-1])

	var q [3]int = [3]int{1, 2, 3}
	for i, v := range a {
		fmt.Printf("%d %d\n", i, v)
	}

	for _, v := range a {
		fmt.Printf("%d\n", v)
	}

	for i, v := range q {
		fmt.Printf("%d|%d\n", i, v)
	}
	s := "hello"
	s += " word"
	fmt.Println(s)

	type Currency int

	const (
		USD Currency = 1 << iota
		EUR
		GBP
		RMB
	)

	symbol := [...]string{USD: "$", EUR: "€", GBP: "￡", RMB: "￥"}

	fmt.Println(USD, symbol[USD])
	fmt.Println(RMB)

	r := [...]int{1: 23}
	for _, v := range r {
		fmt.Println(v)
	}
	println("================")
	fmt.Printf("%d", &r[0])
	fmt.Printf("\n========%d,%d,%d,%d\n", USD, EUR, GBP, RMB)
	slice := make([]StruTest, 0)
	stru := StruTest{name: "zzz"}
	mapp := map[StruTest]int{}
	mapp[stru] = 1
	slice = append(slice, stru)
	fmt.Println(slice[0])
	add := &stru
	fmt.Println(&add)
	fmt.Println(mapp[slice[0]])
	fmt.Println(mapp[stru])
	stru1 := StruTest{name: "zzz"}
	fmt.Println(mapp[stru1])
	directedEdges := make(map[int]int, 0)
	fmt.Println(len(directedEdges))
	directedEdges[0] = 1
	fmt.Println(len(directedEdges))
	directedEdges[2] = 1
	fmt.Println(len(directedEdges))
	directedEdges[3] = 1
	fmt.Println(len(directedEdges))
	directedEdges[4] = 1
	fmt.Println(len(directedEdges))
	fmt.Println(directedEdges[4])

	arr := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(arr[6:])
	mmp := make(map[StruTest]StruTest, 0)
	str := StruTest{name: "d"}
	if mmp[str] == str {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
	//tail := copy()arr[0:]
	////mid:=append(arr[0:0],999);
	//tmp := make([]int, len(tail)+1)
	//arr = append(append(arr[0:0],999),tail...)
	//fmt.Println(arr)
	tmp := make(map[int]struct{})
	fmt.Println(len(tmp))
}

type StruTest struct {
	name string
}
