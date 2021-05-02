package main

import (
	"fmt"
	"time"
)

func main21() {
	chanOne := make(chan int, 1)
	chanTwo := make(chan int, 1)
	go readDate(chanOne)
	fmt.Println(cap(chanTwo))
	fmt.Println("start")
	chanOne <- 1
	chanOne <- 1
	fmt.Println("=====chan1")
	chanTwo <- 2
	fmt.Println("======chan2")
	var count int
	fmt.Println("startSleep")
	time.Sleep(time.Second * 2)
	for {
		select {
		case v := <-chanOne:
			count += v
			fmt.Println("do1")
		case v := <-chanTwo:
			count += v
			fmt.Println("do2")
		}

		if count >= 2 {
			break
		}
	}
	fmt.Println("over")
}

func readDate(ch chan int) {
	for {
		select {
		case v := <-ch:
			fmt.Println(v)
		}
	}
}
