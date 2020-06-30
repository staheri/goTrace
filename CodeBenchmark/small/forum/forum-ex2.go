package main

import (
	"fmt"
	"sync"
)

var sum = 0
var wg sync.WaitGroup

func main() {
	wg.Add(2)
	ch := make(chan int)
	ch <- sum

	go func() {
		sum = <-ch
		sum++
		ch <- sum
		defer wg.Done()
	}()

	go func() {
		sum = <-ch
		sum++
		ch <- sum
		defer wg.Done()
	}()

	wg.Wait()
	i := <-ch
	close(ch)
	fmt.Println(i)
}
