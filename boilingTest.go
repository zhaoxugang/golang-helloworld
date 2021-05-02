/**
声明
*/

package main

import "fmt"

const boilingF = 212.0

func main6() {
	var s int
	var p *bool = nil
	b := false
	p = &b
	str := new(string)
	*str = "ninja"
	*p = true
	println(*p)
	println(*p)
	println(s)
	println(*str)
	var f = boilingF
	var c = (f - 32) * 5 / 9
	fmt.Printf("boiling point = %gF or %gC\n", f, c)
}
