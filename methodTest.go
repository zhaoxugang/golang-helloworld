package main

import (
	"fmt"
	"image/color"
	"math"
)

type Point struct {
	X, Y float64
}
type ColoredPoint struct {
	Point
	Color color.RGBA
}

func main14() {
	p := Point{1, 2}
	q := Point{14, 5}
	println(&p)
	p.Distance(q)
	(&p).Distance(q)
	fmt.Println(p)
	cp := ColoredPoint{Point{1, 2}, color.RGBA{255, 255, 0, 255}}
	fmt.Printf("%cp=%d\n", cp.Point.X)
	distance := Point.Distance
	distance(p, q) //方法表达式
}

// func Distance(p, q Point) float64 {
// 	return math.Hypot(q.X-p.Y, q.Y-p.Y)
// }

func (p Point) Distance(q Point) float64 {
	println(&p)
	p.X = 0
	fmt.Println(p)
	return math.Hypot(q.X-p.Y, q.Y-p.Y)
}
