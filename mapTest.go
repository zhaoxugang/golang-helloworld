package main

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
)

type People struct {
	name string
	age  int32
}

var container2FuncNums cmap.ConcurrentMap

func main18() {
	defer func() {
		fmt.Printf("invoke recover begin\n")
		p := recover()
		if p != nil {
			fmt.Printf("AcquireContainer error，%v", p)
		}
	}()
	fmt.Printf("value=%s\n", runtime.Version)
	fmt.Println(runtime.Version())
	mylist := list.New()
	mylist.PushBack("dasd")
	mylist.PushBack(runtime.Version)

	for element := mylist.Front(); element != nil; element = element.Next() {
		if "dasd"+"1" == element.Value {
			fmt.Printf("value=%s\n", element.Value)
		}
	}
	container2FuncNums = cmap.New()
	ages := make(map[string]int)
	ages1 := map[string]int{
		"alice":   31,
		"charlie": 34,
	}
	defer func() {
		fmt.Printf("invoke recover begin11\n")
		p := recover()
		if p != nil {
			fmt.Printf("AcquireContainer error，%v", p)
		}
	}()
	ages2 := map[string]int{}
	ages2["hello"] = 1
	ages2["world"] = 2
	println(ages)
	println(ages1)
	println(ages2)
	println("asd" + "asd")
	_, ok := ages2["h1ello"]
	println(ok)
	if _, ok := ages2["hello1"]; !ok {
		println("不存在")
	}

	peops := make(map[string]People)
	p1 := peops["赵旭刚"]
	p1.name = "赵旭刚"
	p1.age = 26
	peops["赵旭刚"] = p1
	fmt.Printf("People info = %v\n", peops["赵旭刚"])
	p1 = peops["赵旭刚"]
	p1.age = 26
	fmt.Printf("People info = %v\n", peops)

	map2 := make(map[string]AtomicInt)
	obj, _ := map2["aa"]
	obj.increment()
	map2["aa"] = obj
	fmt.Printf("result=%v", map2)
	obj.increment()
	map2["aa"] = obj
	fmt.Printf("result=%v", map2)
	container2FuncNums.Set("hello", AtomicInt{
		count: 1,
	})
	incrementContainerRequest("hello")
	obj11, _ := container2FuncNums.Get("hello")
	num := obj11.(AtomicInt)
	fmt.Printf("num2=%d\n", num.count)
	incrementContainerRequest("hello")
	fmt.Printf("num3=%d\n", num.count)

}

type AtomicInt struct {
	sync.Mutex
	count int32
}

func (a *AtomicInt) increment() {
	a.Lock()
	a.count += 1
	a.Unlock()
}

func (a *AtomicInt) decrement() {
	a.Lock()
	if a.count > 0 {
		a.count -= 1
	} else {
		a.count = 0
	}
	a.Unlock()
}

func incrementContainerRequest(containerId string) {
	obj, _ := container2FuncNums.Get(containerId)
	num := obj.(AtomicInt)
	num.increment()
	container2FuncNums.Set(containerId, num)
}
