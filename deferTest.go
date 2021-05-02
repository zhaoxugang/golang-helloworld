package main

import (
	"fmt"
	"log"
	"time"
)

func main13() {
	BigSlowOperation()
}

func BigSlowOperation() {
	defer trace("bigSlowOperation")()
	time.Sleep(1 * time.Second)
	DoSth()
}

func DoSth() {
	fmt.Printf("err is coming\n")
	panic("error is catched")
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s\n", msg)
	return func() {
		log.Printf("exit %s(%s)\n", msg, time.Since(start))
		fmt.Printf("recover:%s\n", recover())
	}
}
