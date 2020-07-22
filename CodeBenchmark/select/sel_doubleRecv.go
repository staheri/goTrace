package main

import (
	"fmt"
	"time"
)

func server1(ch chan int) {
	for{
		time.Sleep(2 * time.Millisecond)
		ch <- 1
	}

}
func server2(ch chan int) {
	for{
		time.Sleep(2 * time.Millisecond)
		ch <- 2
	}
}
func main() {
	output1 := make(chan int)
	output2 := make(chan int)
	go server1(output1)
	go server2(output2)
	for {
		select {
		case s1 := <-output1:
			fmt.Println(s1)
		case s2 := <-output2:
			fmt.Println(s2)
		case <-time.After(4 * time.Millisecond):
			return
		}
	}
}
