package main

import (
	"fmt"
  _"time"
)


type Data struct {
	key  int
}

func main() {
	reads := make(chan *Data)
	//done := make(chan int)
	t := Data{0}

	go func() {
		reads <- &t
    reads <- &t
	}()
	go func() {
		t2 := <- reads
		t2.key = 4
		//done <- 0
	}()
	go func() {
		t3 := <- reads
		t3.key = 3
		//done <- 0
	}()
	//<- done
	//<- done
  //time.Sleep(1*time.Second)
	fmt.Println(t.key)

}
