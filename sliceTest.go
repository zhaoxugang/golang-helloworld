package main

import "fmt"

func main8() {
	months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	println(months[0])

	for _, v := range months[2:5] {
		println(v)
	}
	m1 := months[3:4]
	fmt.Printf("==>%d\n", cap(m1))
	m2 := m1[2:9]
	fmt.Printf("<==%d,%d\n", cap(m1), cap(m2))
	fmt.Println(m2)
	println(len(months[3:6]))
	println(cap(months))
	cm := months[3:5]
	cm[1] = 90
	fmt.Printf("%d\n", &cm)
	fmt.Printf("%d\n", &months)
	n := &months[3]
	m := &cm[0]
	println(m)
	println(n)
	println("==============")
	var arr []int = nil
	println(len(arr))
	var slice []int = []int{}
	println(len(slice))
	slice[999] = 12
}
