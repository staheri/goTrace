package main

import (
	"fmt"
	"time"
)

var count = 3

func server1(ch chan int, done chan struct{}) {
	i := 0
	for{
		time.Sleep(3 * time.Microsecond)
		ch <- 1
		i++
		if i > count{
			done <- struct{}{}
			//fmt.Println("DONE1")
			//time.Sleep(1 * time.Millisecond)
			}
	}

}
func server2(ch chan int, done chan struct{}) {
	i := 0
	for{
		time.Sleep(3 * time.Microsecond)
		ch <- 2
		i++
		if i > count{
			done <- struct{}{}
			//fmt.Println("DONE2")
			//time.Sleep(1 * time.Millisecond)
			}
	}
}
func main() {
	output1 := make(chan int)
	output2 := make(chan int)
	done := make(chan struct{})
	go server1(output1,done)
	go server2(output2,done)
forloop:
	for {
		select {
		case s1 := <-output1:
			fmt.Println(s1)
			//s1 = 0
			//continue
			//_ = s1
		case s2 := <-output2:
			fmt.Println(s2)
			//s2 = 0
			//continue
			//_ = s2
		//case <-time.After(10 * time.Microsecond):
		case <- done:
			//time.Sleep(1 * time.Millisecond)
			break forloop
		}
	}
}
