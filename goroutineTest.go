package main

import (
	"fmt"
	"time"
)

func main15() {
	ch := make(chan int64, 100)
	// defer close(ch)
	// go test1(ch)
	// for {
	// 	fmt.Printf("chan len = %d\n", len(ch))
	// 	ch <- time.Now().UnixNano()
	// }
	go produce(ch)
	go consumer(ch)
	for {
		time.Sleep(1000)
	}
}

func test1(ch chan int64) {
	for i := 0; i <= 10; i++ {
		x := <-ch
		fmt.Printf("%d -> [%d]\n", x, i)
	}
}
func produce(ch chan int64) {
	for i := 0; i < 200; i++ {
		ch <- time.Now().UnixNano()
	}
}

func consumer(ch chan int64) {
	for {
		x := <-ch
		fmt.Printf("收到消息[%d]\n", x)
	}
}
