package main

import (
	"fmt"
	"github.com/google/btree"
	selfbtree "helloword/btree"
	"net"
	"sync"
	"time"
)

func main() {
	//b := B{}
	//var a A
	//a = b
	//a.Say()
	//c := a.(A)
	//c.Say()
	//
	//arr := [3]int{1, 2, 3}
	//arr1 := arr[0:0]
	//arr2 := arr[:1]
	//fmt.Print(arr1)
	//fmt.Print(arr2)

	//pop := Peop{}
	//fmt.Print(pop.child)
	//pop.child = append(pop.child, 9)
	//fmt.Print(pop.child)
	//aa := []int{1,2}
	//aa = append(aa, 1)
	//fmt.Println(aa)
	//peop := &Peop{}
	//peop.b = &B{}
	//peop.c = &C{}
	//tmp := peop
	//tmp.b.age = 8
	//fmt.Println("========%v", peop.b.age)
	//fmt.Println("========%v", tmp.b.age)
	//var p1 Peop
	//fmt.Println(p1.b)
	//var str string
	//fmt.Println(str == "")
	//m := 11
	//m &= ^1
	//var p interface{} = nil
	//switch p.(type){
	//case *Peop:
	//	println("PPP")
	//default:
	//	println("gggg")
	//}
	//
	////var err error
	////err = nil
	//_, err := test()
	//println(err)
	//p3 := Peop{}
	//var p4 *Peop
	//p4 = &p3
	//println(p4)
	//buf := getBitStr(-21172)
	//println(fmt.Sprintf("b'%s'", buf))
	//var b uint64
	//b = 0xffffffffffffffff
	//println(b)
	//println(fmt.Sprintf("%s", getBitStr(uint64(0xffffffffffffffff))))

	//dd := Datum{}
	//println(unsafe.Sizeof(dd))
	//println(unsafe.Alignof(dd))
	////println(unsafe.Offsetof(dd))
	//
	//d := 0b11101010111
	//p := d
	//fmt.Printf("%b\n", p)
	//for i := 0; i< 10;i++{
	//	p = (p - 1) & d;
	//	fmt.Printf("%b\n", p)
	//}

	//ch := make(chan int, 3072)
	//go func(ch chan int){
	//	for{
	//		val := <-ch
	//		fmt.Printf("val:%d\n", val)
	//	}
	//}(ch)
	//tick := time.NewTicker(1 * time.Second)
	//for i := 0; i < 20; i++{
	//	select {
	//	case ch<- i:
	//	case <- tick.C:
	//		fmt.Printf("%d:Case<-tick.C\n", i)
	//	}
	//	time.Sleep(200 * time.Millisecond)
	//}
	//close(ch)
	//tick.Stop()

	//var dd *Datum = nil
	//var addr unsafe.Pointer = unsafe.Pointer(dd)
	//fmt.Println(atomic.LoadPointer(&addr) == unsafe.Pointer(nil))
	//fmt.Println(atomic.LoadPointer(&addr))
	//var s pppInt
	//var d Datum
	//d = Datum{
	//	i: 1,
	//}
	//s = &d
	//fmt.Println(&s)

	defer func() {
		fmt.Println("Step2")
	}()
	defer func() {
		fmt.Println("Step1")
	}()

	//tick := time.NewTicker(1 * time.Second)
	now := time.Now().Format("2006-01-02 03:04:05")
	fmt.Println(now)
	//for {
	//	select{
	//	case <- tick.C:
	//		fmt.Println("laile")
	//	}
	//}

	arrs := make([]int, 10)
	arrs[0] = 0
	arrs[1] = 1
	arrs[2] = 2
	arrs[3] = 3
	arrs[4] = 4
	//arrs = arrs[0:5]
	//fmt.Println(arrs[0:3])
	//fmt.Println(arrs[3:5])
	//b := append(arrs[0:3], append([]int{7}, arrs[3:5]...)...)
	//fmt.Println(b)
	//fmt.Println(arrs)
	//c := append(arrs[0:3], []int{7,5,2,2,2,1}...)
	//fmt.Println(c)
	//fmt.Println(arrs)
	//c := arrs[0:2]
	//c = append(c, 9)
	//fmt.Println(c)
	//fmt.Println(arrs)

	selfbtree.BtreeHello()
	selfBtree, err := selfbtree.NewBtree("/Users/zhaoxugang/go/src/golang-helloworld/zdb.idx",
		//3072, false)
		3072, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	var start = 0
	var testNum = 10000000
	//for i := start; i < testNum; i++ {
	//	//n := rand.Int() & 0xffffffff
	//	n := i
	//	if i%50000 == 0 {
	//		fmt.Printf("insert into btree %d, key=%d\n", i, n)
	//	}
	//	//if i == 193 {
	//	//fmt.Println("===")
	//	//}
	//	if i == 3072 {
	//		fmt.Println("ad")
	//	}
	//	key := selfbtree.Encode(n)
	//	value := selfbtree.Encode(n)
	//	ok, _ := selfBtree.Insert(key, value)
	//	if !ok {
	//		fmt.Println("插入失败，key已存在")
	//	}
	//	//if ok {
	//	//	fmt.Println("马蛋！！怎么插入成功了？！")
	//	//}
	//}
	//if err != nil {
	//	fmt.Println(err)
	//}
	for i := start; i < testNum; i++ {
		if i%50000 == 0 {
			fmt.Printf("select from btree %d\n", i)
		}
		//n := randomKeys[i] & 0xffffffff
		n := i
		if i == 2045 {
			fmt.Println("21")
		}
		key := selfbtree.Encode(n)
		value := selfbtree.Encode(n)
		target, ok, err := selfBtree.GET(key)
		if err != nil {
			fmt.Println(err)
		}
		if !ok || !value.Equal(target) {
			fmt.Printf("操蛋,%d\n", i)
		}
	}
	selfBtree.Flush()
	fmt.Println("OVER")
	//fmt.Println("===============别人的BTree====================")
	//bt := btree.New(90)
	//v := btree.Int(2)
	//bt.ReplaceOrInsert(v)
	//fmt.Println(bt.Get(v))
	//var wg *sync.WaitGroup = &sync.WaitGroup{}
	//wg.Add(1)
	//
	//go insertIntoBtree(0, 1000, bt, wg)
	////go insertIntoBtree(1000,2000, bt, wg)
	////go insertIntoBtree(2000,3000, bt, wg)
	////go insertIntoBtree(3000,4000, bt, wg)
	////go insertIntoBtree(4000,5000, bt, wg)
	////go insertIntoBtree(5000,6000, bt, wg)
	//
	//fmt.Println("dd")
	//wg.Wait()
	//fmt.Println("Over")
	//for i := 0; i < 1000; i++{
	//	v := bt.Get(btree.Int(i))
	//	if v != btree.Int(i){
	//		fmt.Println(v)
	//		fmt.Println(i)
	//	}
	//}
	//fmt.Println("START")

	//um, _ := strconv.ParseInt(strconv.Itoa(777), 8, 0)
	//fmt.Println(os.FileMode(777), 777)
	//fmt.Println(os.FileMode(0777), 0777)
	//fmt.Println(os.FileMode(os.ModeDir))
	//fmt.Println(os.FileMode(os.ModeAppend))
	//fmt.Println(os.FileMode(os.ModeExclusive))
	//fmt.Println(os.FileMode(os.ModeTemporary))
	//fmt.Println(os.FileMode(os.ModeSymlink))
	//fmt.Println(os.FileMode(os.ModeDevice))
	//fmt.Println(os.FileMode(os.ModeNamedPipe))
	//fmt.Println(os.FileMode(os.ModeSocket))
	//fmt.Println(os.FileMode(os.ModeSetuid))
	//fmt.Println(os.FileMode(os.ModeSetgid))
	//fmt.Println(os.FileMode(os.ModeCharDevice))
	//fmt.Println(os.FileMode(os.ModeSticky))
	//fmt.Println(os.FileMode(os.ModeIrregular))
}

func insertIntoBtree(start int, end int, bt *btree.BTree, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end; i++ {
		bt.ReplaceOrInsert(btree.Int(i))
	}
}

type pppInt interface {
	saySth()
}

type Datum struct {
	k         byte        // datum kind.
	decimal   uint16      // decimal can hold uint16 values.
	length    uint32      // length can hold uint32 values.
	i         int64       // i can hold int64 uint64 float64 values.
	collation string      // collation hold the collation information for string value.
	b         []byte      // b can hold string or []byte values.
	x         interface{} // x hold all other types.
}

func (d *Datum) saySth() {

}

func getBitStr(n uint64) [64]byte {
	bs := [64]byte{}
	idx := 63
	for n != 0 {
		bs[idx] = byte(n & 1)
		n >>= 1
	}
	return bs
}

func test() (*Peop, error) {
	return nil, &net.DNSError{}
}

type Peop struct {
	b *B
	c *C
}

type A interface {
	Say()
	Work()
}

type B struct {
	age int16
}

type C struct {
	name string
}

func (b B) Say() {
	fmt.Println("say Hello")
}

func (b B) Work() {
	fmt.Println("say bbb")
}
